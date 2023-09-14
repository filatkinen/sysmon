package main

import (
	"flag"
	"log"
	"net"

	"github.com/filatkinen/sysmon/internal/client"
	"github.com/filatkinen/sysmon/internal/model"
)

func main() {
	var (
		everyM   = flag.Int("M", 5, "average period")
		averageN = flag.Int("N", 15, "query period")
		port     = flag.String("port", "50051", "server's port")
		address  = flag.String("address", "localhost", "server's address")
	)

	flag.Parse()

	c, err := NewClientView()
	if err != nil {
		log.Fatalf("error starting cui client: %s", err)
	}

	closeChan := make(chan struct{})

	// Starting cui interface
	go func() {
		err = c.Start()
		if err != nil {
			log.Println(err)
		}
		closeChan <- struct{}{}
	}()

	defer func() {
		c.Stop()
	}()

	// Starting sysmon GRPC client
	cl := client.NewClient(net.JoinHostPort(*address, *port), *everyM, *averageN)
	err = cl.Start()
	if err != nil {
		log.Printf("error starting client: %s", err)
		return
	}
	defer func() {
		err := cl.Close()
		if err != nil {
			log.Printf("got error while closing client:%s", err)
		}
	}()

	// Starting process getting data from GRPC server
	go func() {
		err := cl.GetData(func(data []model.DataToClientStamp) {
			c.GetData(data)
		})
		if err != nil {
			log.Printf("got error while receiving data from server: %s", err)
		}
		closeChan <- struct{}{}
	}()

	<-closeChan
}
