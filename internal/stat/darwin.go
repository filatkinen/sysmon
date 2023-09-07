//go:build darwin

package stat

import "github.com/filatkinen/sysmon/internal/model"

func loadAvg() (model.DataLoadAvg, error) {
	return model.DataLoadAvg{}, ErrNotImplemented
}

func cpuAvgStats() (model.DataCpuAvgStats, error) {
	return model.DataCpuAvgStats{}, ErrNotImplemented
}

func disksLoad() ([]model.DataDisksLoad, error) {
	return []model.DataDisksLoad{}, ErrNotImplemented
}

func disksUsage() ([]model.DataDisksUsage, error) {
	return []model.DataDisksUsage{}, ErrNotImplemented
}

func networkListen() ([]model.DataNetworkListen, error) {
	return []model.DataNetworkListen{}, ErrNotImplemented
}

func networkStates() ([]model.DataNetworkStates, error) {
	return []model.DataNetworkStates{}, ErrNotImplemented
}

func topNetworkProto() ([]model.DataTopNetworkProto, error) {
	return []model.DataTopNetworkProto{}, ErrNotImplemented
}

func topNetworkTraffic() ([]model.DataTopNetworkTraffic, error) {
	return []model.DataTopNetworkTraffic{}, ErrNotImplemented
}
