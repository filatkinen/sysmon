package service_test

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/filatkinen/sysmon/internal/config"
	"github.com/filatkinen/sysmon/internal/model"
	"github.com/filatkinen/sysmon/internal/service"
	"github.com/stretchr/testify/require"
)

func TestServiceCoreLogic(t *testing.T) {
	conf := config.ServiceConfig{
		LA:             true,
		AvgCPU:         false,
		DisksLoad:      false,
		DisksUse:       false,
		NetworkTop:     false,
		NetworkStat:    false,
		ScrapeInterval: time.Second,
		CleanInterval:  time.Minute,
		Depth:          time.Minute * 2,
		Port:           "50052",
		Address:        "0.0.0.0",
	}

	var strangeLock sync.Mutex // added lock due race detection. Though it is not necessary

	srv, err := service.NewService(conf, &MockStat{})
	if err != nil {
		log.Fatalf("error creating service sysmon  %v", err)
	}

	log.Println("Starting test with duration = 6,9 sec (we are expecting 6+1 scrape intervals)")
	ticker := time.NewTicker(time.Second*6 + time.Millisecond*900)
	signalFailedStart := make(chan struct{})

	var errStart, errStop error
	go func() {
		strangeLock.Lock()

		if errStart = srv.Start(); errStart != nil {
			log.Println("failed to start service: " + err.Error())
		}
		strangeLock.Unlock()
		signalFailedStart <- struct{}{}
	}()

	select {
	case <-signalFailedStart:
	case <-ticker.C:
		t.Run("Average time more then having data", func(t *testing.T) {
			_, ready := srv.CountDataClient(20)
			require.False(t, ready)
		})

		t.Run("Average time is equal 2 sec", func(t *testing.T) {
			data, ready := srv.CountDataClient(2)
			require.True(t, ready)
			require.Equal(t, data.Data[0].ElMap["loadavg"][0].NumberField, 6.0, "(4+8)/2=6")
			require.Equal(t, data.Data[0].ElMap["loadavg"][1].NumberField, 12.0, "(8+16)/2=12")
			require.Equal(t, data.Data[0].ElMap["loadavg"][2].NumberField, 18.0, "(12+24)/2=18")
		})

		t.Run("Average time is equal 6 sec", func(t *testing.T) {
			data, ready := srv.CountDataClient(6)
			require.True(t, ready)
			require.Equal(t, data.Data[0].ElMap["loadavg"][0].NumberField, 6.0, "(4+8)*3/6=6")
			require.Equal(t, data.Data[0].ElMap["loadavg"][1].NumberField, 12.0, "(8+16)*3/6=12")
			require.Equal(t, data.Data[0].ElMap["loadavg"][2].NumberField, 18.0, "(12+24)*3/6=18")
		})

		if errStop = srv.Stop(); errStop != nil {
			log.Println("failed to stop service: " + err.Error())
		}
	}
	log.Printf("Exiting sysmon service\n")

	strangeLock.Lock()
	require.NoError(t, errStart)
	require.NoError(t, errStop)
	strangeLock.Unlock()
}

type MockStat struct {
	counter int // odd and even values we are going to send different values
}

func (s *MockStat) LoadAvg() (model.ElMapType, error) {
	m := make(model.ElMapType, 1)
	line := make([]model.Element, 0, 3)
	var el model.Element

	// counter even la1=4 la2=8 la3=12
	// counter odd la1=8 la2=16 la3=24
	for i := 0; i < 3; i++ {
		el.NumberField = float64(i+1) * 4 * float64(s.counter%2+1)
		el.CountAble = true
		el.DecimalField = 0
		line = append(line, el)
	}
	m["loadavg"] = line
	s.counter++
	return m, nil
}

func (s *MockStat) CPUAvgStats() (model.ElMapType, error) {
	return nil, nil
}

func (s *MockStat) DisksLoad() (model.ElMapType, error) {
	return nil, nil
}

func (s *MockStat) DisksUsage() (model.ElMapType, error) {
	return nil, nil
}

func (s *MockStat) NetworkListen() (model.ElMapType, error) {
	return nil, nil
}

func (s *MockStat) NetworkStates() (model.ElMapType, error) {
	return nil, nil
}

func (s *MockStat) TopNetworkProto() (model.ElMapType, error) {
	return nil, nil
}

func (s *MockStat) TopNetworkTraffic() (model.ElMapType, error) {
	return nil, nil
}

func (s *MockStat) Close() {
}
