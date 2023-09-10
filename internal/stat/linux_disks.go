//go:build linux

package stat

import "github.com/filatkinen/sysmon/internal/model"

func disksLoad() (model.ElMapType, error) {
	return nil, nil
}

func disksUsage() (model.ElMapType, error) {
	disks := []struct {
		disk string
		tps  float64
		kbs  float64
	}{{
		disk: "sda1",
		tps:  10,
		kbs:  20,
	}, {
		disk: "sdb2",
		tps:  30,
		kbs:  40,
	}, {
		disk: "nvm5",
		tps:  50,
		kbs:  60,
	}}
	m := make(model.ElMapType, len(disks))
	for _, v := range disks {
		line := make([]model.Element, 0, 3)

		var disk, tps, kbs model.Element

		disk.StringField = v.disk
		line = append(line, disk)

		tps.NumberField = v.tps
		tps.CountAble = true
		tps.DecimalField = 2
		line = append(line, tps)

		kbs.NumberField = v.kbs
		kbs.CountAble = true
		kbs.DecimalField = 2
		line = append(line, tps)

		m[v.disk] = line
	}

	return m, nil
}
