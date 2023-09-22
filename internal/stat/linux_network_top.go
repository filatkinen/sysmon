//go:build linux

package stat

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type netTopTraffic struct {
	sourceIPPort string
	destIPPort   string
	bytes        int
	proto        string
}

type netTopProto struct {
	bytes int
	proto string
}

var (
	netTopTrafficValue = make(map[string]netTopTraffic)
	netTopProtoValue   = make(map[string]netTopProto)

	netTopTrafficCounter int32
	netTopProtoCounter   int32

	wasStartSubSystem bool

	netTopTrafficLastCheck time.Time

	netTopLock sync.Mutex

	ctxCollect    context.Context
	cancelCollect context.CancelFunc
	wgCollect     sync.WaitGroup
	exitChan      = make(chan struct{})

	errorStarting error
)

func topNetworkStartCollect() error {
	netTopLock.Lock()
	defer netTopLock.Unlock()
	if wasStartSubSystem {
		return errorStarting
	}
	wasStartSubSystem = true

	command, err := exec.LookPath("tcpdump")
	//command, err := exec.LookPath("ping")
	if err != nil {
		log.Printf("module networkTop: can not find command %s\n", err)
		errorStarting = err
		return errorStarting
	}

	ctxCollect, cancelCollect = context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctxCollect, command, "-i", "any", "-nt")
	//cmd := exec.CommandContext(ctxCollect, command, "ya.ru")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("module networkTop: can not make cmd.StdoutPipe() %s\n", err)
		errorStarting = err
		return errorStarting
	}

	wgCollect.Add(1)
	go func() {
		defer wgCollect.Done()
		scanner := bufio.NewScanner(cmdReader)

		for scanner.Scan() {
			collect(scanner.Text())
		}

	}()

	if err := cmd.Start(); err != nil {
		e := cmdReader.Close()
		log.Printf("module networkTop: can not run cmd.Start %s\n", errors.Join(err, e))
		errorStarting = err
		return errorStarting
	}

	wgCollect.Add(1)
	go func() {
		defer wgCollect.Done()
		<-exitChan
		cancelCollect()
	}()

	wgCollect.Add(1)
	go func() {
		defer wgCollect.Done()
		if err := cmd.Wait(); err != nil {
			e := ctxCollect.Err()
			if !errors.Is(e, context.Canceled) {
				log.Printf("module networkTop-: got error cmd.Wait  %s\n", err)
				netTopLock.Lock()
				errorStarting = errors.New(fmt.Sprintf("unable to start %s : %s", command, err))
				netTopLock.Unlock()
			}
		}
	}()

	netTopTrafficLastCheck = time.Now()
	return nil
}

func topNetworkTrafficStop() {
	if ctxCollect != nil {
		cancelCollect()
	}
	exitChan <- struct{}{}
	wgCollect.Wait()
}

func collect(str string) {
	//log.Print(str)
	//now := time.Now()
	//defer log.Println(time.Since(now))
	fields := strings.Fields(str)
	lnf := len(fields)
	if lnf < 7 {
		return
	}
	if !(fields[4] == ">" || fields[4] == "<") {
		return
	}

	protoField := fields[6]
	if strings.LastIndex(protoField, ",") == len(protoField)-1 {
		protoField = protoField[:len(protoField)-1]
	}

	if !(protoField == "Flags" || protoField == "UDP" || protoField == "ICMP" || protoField == "ICMP6") {
		return
	}
	if fields[lnf-1] == "0" {
		return
	}

	b, err := strconv.Atoi(fields[lnf-1])
	if err != nil {
		return
	}

	netTopLock.Lock()
	defer netTopLock.Unlock()

	var traf netTopTraffic
	var proto netTopProto

	if protoField == "Flags" {
		protoField = "TCP"
	}

	src := fields[3]
	dst := fields[5][:len(fields[5])-1]
	if protoField == "TCP" || protoField == "UDP" {
		if idx := strings.LastIndex(src, "."); idx != -1 {
			src = src[:idx] + ":" + src[idx+1:]
		}
		if idx := strings.LastIndex(dst, "."); idx != -1 {
			dst = dst[:idx] + ":" + dst[idx+1:]
		}

	}

	traf.sourceIPPort = src
	traf.destIPPort = dst
	traf.proto = protoField
	traf.bytes = b

	proto.proto = protoField
	proto.bytes = b

	hash := traf.sourceIPPort + traf.destIPPort + traf.proto
	trafVal, ok := netTopTrafficValue[hash]
	if ok {
		trafVal.bytes += b
		netTopTrafficValue[hash] = trafVal
	} else {
		netTopTrafficValue[hash] = traf
	}

	protoVal, ok := netTopProtoValue[proto.proto]
	if ok {
		protoVal.bytes += b
		netTopProtoValue[proto.proto] = protoVal
	} else {
		netTopProtoValue[proto.proto] = proto
	}

	if protoField == "ICMP" {
		log.Print("---------", traf, proto)
	}
}
