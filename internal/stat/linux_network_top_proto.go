//go:build linux

package stat

import (
	"github.com/filatkinen/sysmon/internal/model"
)

func topNetworkProto() (model.ElMapType, error) {
	headersLen := len(model.StampNameHeaders[6].Header)

	err := topNetworkStartCollect()
	if err != nil {
		return returnError(headersLen, err)
	}

	if netTopProtoCounter == 0 {
		netTopProtoCounter++
		return returnZeroSlice(headersLen)
	}

	netTopLock.Lock()
	defer netTopLock.Unlock()

	m := make(model.ElMapType, len(netTopProtoValue))

	for k, v := range netTopProtoValue {
		line := make([]model.Element, 0, headersLen)
		var el model.Element

		el.StringField = v.proto
		line = append(line, el)

		el.StringField = ""
		el.CountAble = true
		el.NumberField = float64(v.bytes)
		el.DecimalField = 0
		line = append(line, el)

		el.StringField = ""
		el.CountAble = true
		el.PercentAble = true
		el.NumberField = float64(v.bytes)
		el.DecimalField = 0
		line = append(line, el)

		m[k] = line
	}

	clear(netTopProtoValue)
	return m, nil
	//return returnZeroSlice(headersLen)
}
