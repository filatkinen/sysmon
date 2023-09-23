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
