package stat

import (
	"errors"
	"github.com/filatkinen/sysmon/internal/model"
	"sync"
	"time"
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

type Stat struct {
	netTopOnceInit    sync.Once
	netTopInitLock    sync.Mutex
	netTopData        sync.Mutex
	netTopLastCheck   time.Time
	netTopIsEnable    bool
	exitChan          sync.Mutex
	netTopTrafficData []netTopTraffic
}

func New() *Stat {
	return &Stat{
		netTopOnceInit: sync.Once{},
		netTopInitLock: sync.Mutex{},
		netTopData:     sync.Mutex{},
		netTopIsEnable: false,
		exitChan:       sync.Mutex{},
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
	return topNetworkProto()
}

func (s *Stat) TopNetworkTraffic() (model.ElMapType, error) {
	return topNetworkTraffic()
}

func (s *Stat) Close() error {
	return nil
}

func (s *Stat) topNetworkStartSubsystem() error {
	return nil
}

func returnError(elCount int, err error) (model.ElMapType, error) {
	m := make(model.ElMapType, 1)
	line := make([]model.Element, elCount)
	line[0].StringField = err.Error()
	m["error"] = line
	return m, err
}
