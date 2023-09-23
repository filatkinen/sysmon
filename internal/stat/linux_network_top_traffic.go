//go:build linux

package stat

import (
	"github.com/filatkinen/sysmon/internal/model"
	"time"
)

func topNetworkTraffic() (model.ElMapType, error) {
	headersLen := len(model.StampNameHeaders[7].Header)

	err := topNetworkStartCollect()
	if err != nil {
		return returnError(headersLen, err)
	}

	if netTopTrafficCounter == 0 {
		netTopTrafficCounter++
		return returnZeroSlice(headersLen)
	}
	netTopLock.Lock()
	defer netTopLock.Unlock()

	now := time.Now()
	seconds := now.Sub(netTopTrafficLastCheck).Seconds()

	m := make(model.ElMapType, len(netTopTrafficValue))
	for k, v := range netTopTrafficValue {
		line := make([]model.Element, 0, headersLen)
		var el model.Element

		el.StringField = v.sourceIPPort
		line = append(line, el)

		el.StringField = v.destIPPort
		line = append(line, el)

		el.StringField = v.proto
		line = append(line, el)

		el.CountAble = true
		el.NumberField = float64(v.bytes) / (seconds)
		el.DecimalField = 0
		line = append(line, el)

		m[k] = line
	}

	netTopTrafficLastCheck = now
	clear(netTopTrafficValue)

	return m, nil
}
