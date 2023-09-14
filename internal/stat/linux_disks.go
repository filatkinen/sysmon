//go:build linux

package stat

import (
	"math/rand"
	"time"

	"github.com/filatkinen/sysmon/internal/model"
)

func disksLoad() (model.ElMapType, error) {
	return nil, nil
}

func disksUsage() (model.ElMapType, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1) //nolint
	n1 := r1.Intn(5) + 1
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
		if n1 == 3 {
			disk.StringField = v.disk + "/New"
		} else {
			disk.StringField = v.disk
		}
		line = append(line, disk)

		tps.NumberField = v.tps * float64(n1)
		tps.CountAble = true
		tps.DecimalField = 2
		line = append(line, tps)

		kbs.NumberField = v.kbs * float64(n1)
		kbs.CountAble = true
		kbs.DecimalField = 2
		line = append(line, kbs)

		if n1 == 3 {
			m[v.disk+"/New"] = line
		} else {
			m[v.disk] = line
		}
	}

	return m, nil
}
