package stat

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/filatkinen/sysmon/internal/model"
	"os/exec"
	"strings"
)

type NetStates struct {
	Status string
	Number float64
}

func networkStates() (model.ElMapType, error) {
	headersLen := len(model.StampNameHeaders[5].Header)

	netsStates, err := networkStatesQuery()
	if err != nil {
		return returnError(headersLen, err)
	}

	m := make(model.ElMapType, len(netsStates))
	for _, v := range netsStates {
		line := make([]model.Element, 0, headersLen)
		var el model.Element
		// "(Status)", "(Number)"
		el.StringField = v.Status
		line = append(line, el)

		el.StringField = ""
		el.CountAble = true
		el.DecimalField = 0
		el.NumberField = v.Number
		line = append(line, el)

		m[v.Status] = line
	}
	return m, nil

}

func networkStatesQuery() ([]NetStates, error) {
	netStatesSlice := make([]NetStates, 0, 16)
	ss, err := exec.LookPath("ss")
	if err != nil {
		return nil, err
	}

	out, err := exec.Command(ss, "-at").Output()

	if err != nil {
		return nil, err
	}

	m := make(map[string]float64)

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			if err != nil {
				return nil, errors.New("got error parsing status using ss command")
			}
		}
		netStatus := fields[0]
		if err != nil {
			return nil, err
		}
		m[netStatus] = m[netStatus] + 1
	}
	for k, v := range m {
		netStatesSlice = append(netStatesSlice, NetStates{
			Status: k,
			Number: v,
		})
	}
	return netStatesSlice, nil
}
