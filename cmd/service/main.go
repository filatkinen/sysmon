package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/filatkinen/sysmon/internal/config"
	"github.com/filatkinen/sysmon/internal/service"
	"github.com/filatkinen/sysmon/internal/stat"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "configs/service.yaml", "Path to configuration file")
	flag.Parse()

	conf, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("error reading config file %v", err)
	}

	srv, err := service.NewService(conf, stat.Stat{})
	if err != nil {
		log.Fatalf("error creating service sysmon  %v", err)
	}

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	signalFailedStart := make(chan struct{})
	go func() {
		if err = srv.Start(); err != nil {
			log.Println("failed to start service: " + err.Error())
		}
		signalFailedStart <- struct{}{}
	}()

	select {
	case <-signalFailedStart:
	case sig := <-exitChan:
		log.Printf("Got exit signal %d. Exiting sysmon service\n", sig)
		if err = srv.Stop(); err != nil {
			log.Println("failed to stop service: " + err.Error())
		}
	}
	log.Printf("Exiting sysmon service\n")
}
