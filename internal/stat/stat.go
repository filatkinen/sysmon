package stat

import (
	"context"
	"errors"
	"github.com/filatkinen/sysmon/internal/model"
	"sync"
)

var ErrNotImplemented = errors.New("not implemented")

type netTopTraffic struct {
	sourceIP   string
	sourcePort string
	destIP     string
	destPort   string
	bytes      uint64
	proto      string
}

type netTopProto struct {
	sourceIP   string
	sourcePort string
	destIP     string
	destPort   string
	bytes      uint64
	proto      string
}

type Stat struct {
	netTopOnceInit sync.Once
	//netTopInitLock       sync.Mutex
	//netTopData           sync.Mutex

	ctxTcpDumpCmd        context.Context
	ctxTcpDumpCancelFunc context.CancelFunc

	netTopIsEnable bool

	exitChan          chan struct{}
	netTopTrafficData []netTopTraffic
	wg                sync.WaitGroup
}

func New() *Stat {
	return &Stat{
		netTopOnceInit: sync.Once{},
		//netTopInitLock: sync.Mutex{},
		//netTopData:     sync.Mutex{},
		netTopIsEnable: false,
		exitChan:       make(chan struct{}),
		wg:             sync.WaitGroup{},
	}
}

func (s *Stat) LoadAvg() (model.ElMapType, error) {
	return loadAvg()
}

func (s *Stat) CPUAvgStats() (model.ElMapType, error) {
	return cpuAvgStats()
}

func (s *Stat) DisksLoad() (model.ElMapType, error) {
	return disksLoad()
}

func (s *Stat) DisksUsage() (model.ElMapType, error) {
	return disksUsage()
}

func (s *Stat) NetworkListen() (model.ElMapType, error) {
	return networkListen()
}

func (s *Stat) NetworkStates() (model.ElMapType, error) {
	return networkStates()
}

func (s *Stat) TopNetworkProto() (model.ElMapType, error) {
	s.netTopOnceInit.Do(s.topNetworkStartSubsystem)
	return topNetworkProto(bool)
}

func (s *Stat) TopNetworkTraffic() (model.ElMapType, error) {
	s.netTopOnceInit.Do(s.topNetworkStartSubsystem)
	return topNetworkTraffic()
}

func (s *Stat) Close() {
	s.exitChan <- struct{}{}
	s.wg.Wait()
}

func returnError(elCount int, err error) (model.ElMapType, error) {
	m := make(model.ElMapType, 1)
	line := make([]model.Element, elCount)
	line[0].StringField = err.Error()
	m["error"] = line
	return m, err
}
