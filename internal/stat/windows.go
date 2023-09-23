//go:build windows

package stat

import "github.com/filatkinen/sysmon/internal/model"

func loadAvg() (model.ElMapType, error) {
	return returnError(0, ErrNotImplemented)
}

func cpuAvgStats() (model.ElMapType, error) {
	return returnError(1, ErrNotImplemented)
}

func disksLoad() (model.ElMapType, error) {
	return returnError(2, ErrNotImplemented)
}

func disksUsage() (model.ElMapType, error) {
	return returnError(3, ErrNotImplemented)
}

func networkListen() (model.ElMapType, error) {
	return returnError(4, ErrNotImplemented)
}

func networkStates() (model.ElMapType, error) {
	return returnError(5, ErrNotImplemented)
}

func topNetworkProto() (model.ElMapType, error) {
	return returnError(6, ErrNotImplemented)
}

func topNetworkTraffic() (model.ElMapType, error) {
	return returnError(7, ErrNotImplemented)
}
