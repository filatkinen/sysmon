//go:build linux

package stat

import (
	"github.com/filatkinen/sysmon/internal/model"
	"os"
	"strconv"
	"strings"
)

func loadAvg() (model.ElMapType, error) {
	m := make(model.ElMapType, 1)
	line := make([]model.Element, 0, 3)
	var el model.Element

	file, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return returnError(cap(line), err)
	}

	fields := strings.Fields(string(file))
	loadAvg1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return returnError(cap(line), err)
	}
	el.NumberField = loadAvg1
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	loadAvg5, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return returnError(cap(line), err)
	}
	el.NumberField = loadAvg5
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	loadAvg15, err := strconv.ParseFloat(fields[2], 64)
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
