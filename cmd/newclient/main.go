package main

import (
	"flag"
	"fmt"
	"github.com/filatkinen/sysmon/internal/client"
	"github.com/filatkinen/sysmon/internal/model"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var everyM = flag.Int("M", 5, "average period")
	var averageN = flag.Int("N", 15, "query period")
	var port = flag.String("port", "50051", "server's port")
	var address = flag.String("address", "localhost", "server's address")

	flag.Parse()

	cl := client.NewClientSysMon(net.JoinHostPort(*address, *port), *everyM, *averageN)
	err := cl.Start()
	if err != nil {
		log.Fatalf("error starting client: %s", err)
	}
	defer func() {
		err := cl.Close()
		if err != nil {
			log.Printf("got error while closing client:%s", err)
		}
	}()

	closeGettingDataChan := make(chan struct{})
	go func() {
		err := cl.GetData(showData)
		if err != nil {
			log.Printf("got error while recieving data from server: %s", err)
		}
		closeGettingDataChan <- struct{}{}
	}()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case sig := <-exitChan:
		log.Printf("got exit signal %d. Exiting sysmon service\n", sig)
	case <-closeGettingDataChan:
	}
}

func showData(m []model.DataToClientStamp) {
	for i := range m {
		fmt.Printf("Parameter:%s\n", m[i].Name)
		for k1 := range m[i].Data {
			for k2 := range m[i].Data[k1] {
				fmt.Printf("\t%s", m[i].Data[k1][k2])
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
