//go:build linux

package stat

import "github.com/filatkinen/sysmon/internal/model"

func cpuAvgStats() (model.ElMapType, error) {
	m := make(model.ElMapType, 1)
	line := make([]model.Element, 0, 3)

	var el model.Element

	el.NumberField = 10.0
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	el.NumberField = 20.0
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	el.NumberField = 30.0
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	m["cpuAvgStats"] = line
	return m, nil
}
