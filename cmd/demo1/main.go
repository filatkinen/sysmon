package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
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
	netstat, err := exec.LookPath("sar")
	if err != nil {
		return
	}

	//cmd := exec.Command(netstat, "-c")
	cmd := exec.Command(netstat, "-n", "DEV", "1", "4")
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
