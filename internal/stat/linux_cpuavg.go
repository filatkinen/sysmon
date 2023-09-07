//go:build linux

package stat

import "github.com/filatkinen/sysmon/internal/model"

func cpuAvgStats() (model.DataCpuAvgStats, error) {
	return model.DataCpuAvgStats{}, nil
}
