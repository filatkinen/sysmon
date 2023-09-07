package internal

import "github.com/filatkinen/sysmon/internal/model"

type StatGetter interface {
	LoadAvg() (model.DataLoadAvg, error)
	CpuAvgStats() (model.DataCpuAvgStats, error)
	DisksLoad() ([]model.DataDisksLoad, error)
	DisksUsage() ([]model.DataDisksUsage, error)
	NetworkListen() ([]model.DataNetworkListen, error)
	NetworkStates() ([]model.DataNetworkStates, error)
	TopNetworkProto() ([]model.DataTopNetworkProto, error)
	TopNetworkTraffic() ([]model.DataTopNetworkTraffic, error)
}
