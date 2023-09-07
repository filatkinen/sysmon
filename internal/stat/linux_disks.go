//go:build linux

package stat

import "github.com/filatkinen/sysmon/internal/model"

func disksLoad() ([]model.DataDisksLoad, error) {
	return []model.DataDisksLoad{}, nil
}

func disksUsage() ([]model.DataDisksUsage, error) {
	return []model.DataDisksUsage{}, nil
}
