package stat

import (
	"errors"
	"github.com/filatkinen/sysmon/internal/model"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type Stat struct {
}

func (s Stat) LoadAvg() (model.DataLoadAvg, error) {
	return loadAvg()
}

func (s Stat) CpuAvgStats() (model.DataCpuAvgStats, error) {
	return cpuAvgStats()
}

func (s Stat) DisksLoad() ([]model.DataDisksLoad, error) {
	return disksLoad()
}

func (s Stat) DisksUsage() ([]model.DataDisksUsage, error) {
	return disksUsage()
}

func (s Stat) NetworkListen() ([]model.DataNetworkListen, error) {
	return networkListen()
}

func (s Stat) NetworkStates() ([]model.DataNetworkStates, error) {
	return networkStates()
}

func (s Stat) TopNetworkProto() ([]model.DataTopNetworkProto, error) {
	return topNetworkProto()
}

func (s Stat) TopNetworkTraffic() ([]model.DataTopNetworkTraffic, error) {
	return topNetworkTraffic()
}
