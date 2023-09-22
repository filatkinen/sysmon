//go:build linux

package stat

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/filatkinen/sysmon/internal/model"
	"os/exec"
	"strconv"
	"strings"
)

type NetListen struct {
	Ptoto        string
	Command      string
	Pid          string
	User         string
	LocalAddress string
	Port         string
}

func networkListen() (model.ElMapType, error) {
	headersLen := len(model.StampNameHeaders[4].Header)

	netsListen, err := networkListenQuery()
	if err != nil {
		return returnError(headersLen, err)
	}

	m := make(model.ElMapType, len(netsListen))
	for _, v := range netsListen {
		line := make([]model.Element, 0, headersLen)
		var el model.Element
		var sb strings.Builder
		// "(Protocol)", "(Command)", "(PID)", "(USER)", "(Local Address)", "(PORT)")
		el.StringField = v.Ptoto
		sb.WriteString(el.StringField)
		line = append(line, el)

		el.StringField = v.Command
		sb.WriteString(el.StringField)
		line = append(line, el)

		el.StringField = v.Pid
		sb.WriteString(el.StringField)
		line = append(line, el)

		el.StringField = v.User
		sb.WriteString(el.StringField)
		line = append(line, el)

		el.StringField = v.LocalAddress
		sb.WriteString(el.StringField)
		line = append(line, el)

		el.StringField = ""
		el.CountAble = true
		el.DecimalField = 0
		el.NumberField, err = strconv.ParseFloat(v.Port, 64)
		if err != nil {
			return returnError(headersLen, err)
		}
		sb.WriteString(v.Port)
		line = append(line, el)

		m[sb.String()] = line
	}
	return m, nil

}

func networkListenQuery() ([]NetListen, error) {
	netListenSlice := make([]NetListen, 0, 16)
	netstat, err := exec.LookPath("netstat")
	if err != nil {
		return nil, err
	}

	out, err := exec.Command(netstat, "-ntlpue").Output()

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()
		netListen, err := networkListenQueryParse(line)
		if err != nil {
			return nil, err
		}
		netListenSlice = append(netListenSlice, netListen)
	}

	return netListenSlice, nil
}

func networkListenQueryParse(in string) (NetListen, error) {
	fields := strings.Fields(in)

	if len(fields) < 8 {
		return NetListen{}, errors.New("couldn't parse netstat because there less than 8 fields")
	}
	netListen := NetListen{}
	// Proto Recv-Q Send-Q Local Address           Foreign Address         State       User       Inode      PID/Program name
	// tcp        0      0 0.0.0.0:9090            0.0.0.0:*               LISTEN      0          24470      2155/docker-proxy: 1
	// tcp        0      0 127.0.0.54:53           0.0.0.0:*               LISTEN      996        34126      717/systemd-resolve
	if fields[5] == "LISTEN" {
		fields = append(fields[:5], fields[6:]...)
	}
	for i, v := range fields {
		switch i {
		case 0:
			netListen.Ptoto = v
		case 3:
			idx := strings.LastIndex(v, ":")
			if idx == -1 {
				return NetListen{}, errors.New("got error parsing local address")
			}
			netListen.Port = v[idx+1:]
			netListen.LocalAddress = v[:idx]
		case 5:
			netListen.User = v
		case 7:
			if v == "-" {
				netListen.Pid = v
				netListen.Command = v
				continue
			}
			idx := strings.Index(v, "/")
			if idx == -1 {
				return NetListen{}, errors.New("got error parsing PID/Program name")
			}
			netListen.Pid = v[:idx]
			netListen.Command = v[idx+1:]
			lenC := len(netListen.Command)
			if lenC > 0 || netListen.Command[lenC-1:] == ":" {
				netListen.Command = netListen.Command[:lenC-1]
			}
		}
	}
	return netListen, nil
}
