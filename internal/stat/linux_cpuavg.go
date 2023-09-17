//go:build linux

package stat

import (
	"bufio"
	"github.com/filatkinen/sysmon/internal/model"
	"os"
	"strconv"
	"strings"
	"time"
)

func cpuAvgStats() (model.ElMapType, error) {
	m := make(model.ElMapType, 1)
	line := make([]model.Element, 0, 3)
	var el model.Element

	total1, values1, err := getCPUStats()
	if err != nil {
		return returnError(cap(line), err)
	}

	time.Sleep(time.Millisecond * 100)

	total2, values2, err := getCPUStats()
	if err != nil {
		return returnError(cap(line), err)
	}

	delta := abs(total2 - total1)

	for i := 0; i < 3; i++ {
		el.NumberField = abs(values1[i]-values2[i]) * 100 / delta
		el.CountAble = true
		el.DecimalField = 2
		el.StringField = ""
		line = append(line, el)
	}

	m["cpuAvgStats"] = line
	return m, nil
}

func getCPUStats() (total float64, stat []float64, err error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0.0, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	fields := strings.Fields(scanner.Text())
	for i := range fields {
		if i == 0 {
			continue
		}
		value, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return 0.0, nil, err
		}
		total += value
		switch {
		case i == 1: // User
			stat = append(stat, value)
		case i == 3: // System
			stat = append(stat, value)
		case i == 4: // Idle
			stat = append(stat, value)
		}
	}
	return total, stat, nil
}

func abs(f float64) float64 {
	if f > 0 {
		return f
	}
	return -1 * f
}
