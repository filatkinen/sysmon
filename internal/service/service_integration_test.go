//go:build integrational

package service_test

import (
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/filatkinen/sysmon/internal/config"
	"github.com/filatkinen/sysmon/internal/model"
	"github.com/filatkinen/sysmon/internal/service"
	"github.com/filatkinen/sysmon/internal/stat"
	"github.com/stretchr/testify/require"
)

func TestServiceIntegration(t *testing.T) {
	conf := config.ServiceConfig{
		LA:             false,
		AvgCPU:         false,
		DisksLoad:      false,
		DisksUse:       false,
		NetworkTop:     false,
		NetworkStat:    true,
		ScrapeInterval: time.Second,
		CleanInterval:  time.Minute,
		Depth:          time.Minute * 2,
		Port:           "50053",
		Address:        "0.0.0.0",
	}

	srv, err := service.NewService(conf, &stat.Stat{})
	if err != nil {
		log.Fatalf("error creating service sysmon  %v", err)
	}

	log.Println("Starting integration test...")

	signalExit := make(chan struct{})
	signalFailedStart := make(chan struct{})
	var errStart, errStop error
	go func() {
		if errStart = srv.Start(); errStart != nil {
			log.Println("failed to start service: " + err.Error())
		}
		signalFailedStart <- struct{}{}
	}()

	ticker := time.NewTicker(time.Second * 4)

	select {
	case <-signalFailedStart:
		close(signalExit)
	case <-ticker.C:
		const addListen = 10
		var wg sync.WaitGroup
		var countListen int32
		dataBefore, ready := srv.CountDataClient(1)
		require.True(t, ready)
		wg.Add(addListen)
		for i := 0; i < addListen; i++ {
			go func(i int) {
				defer wg.Done()
				l, err := net.Listen("tcp", net.JoinHostPort("localhost", strconv.Itoa(57000+i)))
				if err != nil {
					return
				}
				atomic.AddInt32(&countListen, 1)
				defer l.Close()
				<-signalExit
			}(i)
		}

		time.Sleep(time.Second * 2)

		t.Run("Compare Listen states before and after", func(t *testing.T) {
			dataAfter, ready := srv.CountDataClient(1)
			require.True(t, ready)
			require.Equal(t, getListen(dataBefore)+atomic.LoadInt32(&countListen), getListen(dataAfter))
		})
		if errStop = srv.Stop(); errStop != nil {
			log.Println("failed to stop service: " + err.Error())
		}
		close(signalExit)
		wg.Wait()
	}
	log.Printf("Exiting sysmon service\n")

	require.NoError(t, errStart)
	require.NoError(t, errStop)
}

func getListen(d *model.StampsData) int32 {
	for i, _ := range d.Data {
		if d.Data[i].IdxStampNameHeaders == 5 {
			return int32(d.Data[i].ElMap["LISTEN"][1].NumberField)
		}
	}
	return 0
}
