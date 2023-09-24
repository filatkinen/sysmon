//go:build linux

package stat

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/filatkinen/sysmon/internal/model"
)

var errorSkipLine = errors.New("skip line")

type DiskLoad struct {
	Device string
	Tps    float64
	KbRW   float64
}

func disksLoad() (model.ElMapType, error) {
	headersLen := len(model.StampNameHeaders[2].Header)

	disksload, err := disksLoadQuery()
	if err != nil {
		return returnError(headersLen, err)
	}

	m := make(model.ElMapType, len(disksload))
	for _, v := range disksload {
		line := make([]model.Element, 0, headersLen)
		var el model.Element

		// "(Device)", "(Tps)", "(Kbps)"
		el.StringField = v.Device
		line = append(line, el)

		el.NumberField = v.Tps
		el.CountAble = true
		el.DecimalField = 2
		line = append(line, el)

		el.NumberField = v.KbRW
		el.CountAble = true
		el.DecimalField = 2
		line = append(line, el)

		m[v.Device] = line
	}
	return m, nil
}

func disksLoadQuery() ([]DiskLoad, error) {
	iostat, err := exec.LookPath("iostat")
	if err != nil {
		return nil, err
	}

	out, err := exec.Command(iostat, "-yd", "1", "1").Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Scan()
	scanner.Scan()
	scanner.Scan()
	scanner.Scan()
	disksload := make([]DiskLoad, 0, 2)
	for scanner.Scan() {
		line := scanner.Text()
		diskload, err := disksLoadQueryParse(line)
		if err != nil {
			if errors.Is(err, errorSkipLine) {
				continue
			}
			return nil, err
		}
		disksload = append(disksload, diskload)
	}
	return disksload, nil
}

func disksLoadQueryParse(load string) (DiskLoad, error) {
	fields := strings.Fields(load)
	if len(fields) == 0 {
		return DiskLoad{}, errorSkipLine
	}
	if strings.Contains(fields[0], "loop") {
		return DiskLoad{}, errorSkipLine
	}
	var diskLoad DiskLoad
	diskLoad.Device = fields[0]
	tps, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return DiskLoad{}, err
	}
	diskLoad.Tps = tps
	kbr, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return DiskLoad{}, err
	}
	kbw, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return DiskLoad{}, err
	}
	diskLoad.KbRW = kbr + kbw

	return diskLoad, nil
}
