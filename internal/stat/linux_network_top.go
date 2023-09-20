//go:build linux

package stat

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func (s *Stat) topNetworkStartSubsystem() {

	tcpdump, err := exec.LookPath("ping")
	if err != nil {
		log.Printf("can not find tcpdump %s\n", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.ctxTcpDumpCmd = ctx
	s.ctxTcpDumpCancelFunc = cancel

	cmd := exec.CommandContext(ctx, tcpdump, "ya.ru")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("can not make  cmd.StdoutPipe() %s\n", err)
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		scanner := bufio.NewScanner(cmdReader)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

	}()

	if err := cmd.Start(); err != nil {
		e := cmdReader.Close()
		log.Printf("can not Start cmd %s\n", errors.Join(err, e))
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		<-s.exitChan
		s.ctxTcpDumpCancelFunc()
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := cmd.Wait(); err != nil {
			log.Printf("got error cmd.Wait stat module: %s\n", err)
		}
	}()
	s.netTopIsEnable = true
}
