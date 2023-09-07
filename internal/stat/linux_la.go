//go:build linux

package stat

import "github.com/filatkinen/sysmon/internal/model"

func loadAvg() (model.DataLoadAvg, error) {
	return model.DataLoadAvg{}, nil
}
