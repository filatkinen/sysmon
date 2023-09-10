//go:build linux

package stat

import "github.com/filatkinen/sysmon/internal/model"

func loadAvg() (model.ElMapType, error) {
	m := make(model.ElMapType, 1)
	line := make([]model.Element, 0, 3)

	var el model.Element
	el.NumberField = 2.0
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	el.NumberField = 4.0
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	el.NumberField = 6.0
	el.CountAble = true
	el.DecimalField = 2
	el.StringField = ""
	line = append(line, el)

	m["loadavg"] = line
	return m, nil
}
