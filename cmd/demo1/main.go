package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func main() {
	//l, err := net.Interfaces()
	//if err != nil {
	//	panic(err)
	//
	//}
	//for _, f := range l {
	//	fmt.Println(f.Name)
	//}
	//netstat, err := exec.LookPath("netstat")
	netstat, err := exec.LookPath("ping")
	if err != nil {
		return
	}

	//cmd := exec.Command(netstat, "-c")
	//cmd := exec.Command(netstat, "ya.ru")
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, netstat, "ya.ru")
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

	go func() {
		time.Sleep(time.Second * 3)
		cancel()
	}()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	//out, err := exec.Command(netstat, "-c").Output()
	//if err != nil {
	//	return
	//
	//	scanner := bufio.NewScanner(bytes.NewReader(out))
	//	for scanner.Scan() {
	//		line := scanner.Text()
	//		fmt.Println(line)
	//	}
	//
	//}
}
