//go:build darwin

package stat

import (
	"errors"
	"github.com/filatkinen/sysmon/internal/model"
	"os/exec"
	"strconv"
	"strings"
)

func loadAvg() (model.ElMapType, error) {
	m := make(model.ElMapType, 1)
	line := make([]model.Element, 0, 3)
	var el model.Element

	out, err := exec.Command(`sysctl`, `-n`, `vm.loadavg`).Output()
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(out))
	if len(fields) < 4 {
		return returnError(cap(line), errors.New("error parsing loadavg"))
	}
	loadAvg1, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return returnError(cap(line), err)
	}
	el.NumberField = loadAvg1
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	loadAvg5, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return returnError(cap(line), err)
	}
	el.NumberField = loadAvg5
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	loadAvg15, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return returnError(cap(line), err)
	}
	el.NumberField = loadAvg15
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	m["loadavg"] = line
	return m, nil
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
