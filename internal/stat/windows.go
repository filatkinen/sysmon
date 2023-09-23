//go:build windows

package stat

import "github.com/filatkinen/sysmon/internal/model"

func loadAvg() (model.ElMapType, error) {
	return returnError(len(model.StampNameHeaders[0].Header), ErrNotImplemented)
}

func cpuAvgStats() (model.ElMapType, error) {
	return returnError(len(model.StampNameHeaders[1].Header), ErrNotImplemented)
}

func disksLoad() (model.ElMapType, error) {
	return returnError(len(model.StampNameHeaders[2].Header), ErrNotImplemented)
}

func disksUsage() (model.ElMapType, error) {
	return returnError(len(model.StampNameHeaders[3].Header), ErrNotImplemented)
}

func networkListen() (model.ElMapType, error) {
	return returnError(len(model.StampNameHeaders[4].Header), ErrNotImplemented)
}

func networkStates() (model.ElMapType, error) {
	return returnError(len(model.StampNameHeaders[5].Header), ErrNotImplemented)
}

func topNetworkProto() (model.ElMapType, error) {
	return returnError(len(model.StampNameHeaders[6].Header), ErrNotImplemented)
}

func topNetworkTraffic() (model.ElMapType, error) {
	return returnError(len(model.StampNameHeaders[7].Header), ErrNotImplemented)
}

func topNetworkTrafficStop() {

}
