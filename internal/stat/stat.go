package stat

import (
	"errors"

	"github.com/filatkinen/sysmon/internal/model"
)

var ErrNotImplemented = errors.New("not implemented")

type Stat struct{}

func (s Stat) LoadAvg() (model.ElMapType, error) {
	return loadAvg()
}

func (s Stat) CPUAvgStats() (model.ElMapType, error) {
	return cpuAvgStats()
}

func (s Stat) DisksLoad() (model.ElMapType, error) {
	return disksLoad()
}

func (s Stat) DisksUsage() (model.ElMapType, error) {
	return disksUsage()
}

func (s Stat) NetworkListen() (model.ElMapType, error) {
	return networkListen()
}

func (s Stat) NetworkStates() (model.ElMapType, error) {
	return networkStates()
}

func (s Stat) TopNetworkProto() (model.ElMapType, error) {
	return topNetworkProto()
}

func (s Stat) TopNetworkTraffic() (model.ElMapType, error) {
	return topNetworkTraffic()
}
