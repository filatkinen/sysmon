//go:build linux

package stat

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

var onceNetworkTop sync.Once
var isNetworkStarted bool

func startNetWorkTop() {
	tcpdump, err := exec.LookPath("ping")
	if err != nil {
		return
	}

	cmd := exec.Command(tcpdump, "ya.ru")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

}
