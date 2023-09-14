package internal

import "github.com/filatkinen/sysmon/internal/model"

type StatGetter interface {
	LoadAvg() (model.ElMapType, error)
	CPUAvgStats() (model.ElMapType, error)
	DisksLoad() (model.ElMapType, error)
	DisksUsage() (model.ElMapType, error)
	NetworkListen() (model.ElMapType, error)
	NetworkStates() (model.ElMapType, error)
	TopNetworkProto() (model.ElMapType, error)
	TopNetworkTraffic() (model.ElMapType, error)
}
