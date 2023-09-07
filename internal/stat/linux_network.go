//go:build linux

package stat

import "github.com/filatkinen/sysmon/internal/model"

func networkListen() ([]model.DataNetworkListen, error) {
	return []model.DataNetworkListen{}, nil
}

func networkStates() ([]model.DataNetworkStates, error) {
	return []model.DataNetworkStates{}, nil
}

func topNetworkProto() ([]model.DataTopNetworkProto, error) {
	return []model.DataTopNetworkProto{}, nil
}

func topNetworkTraffic() ([]model.DataTopNetworkTraffic, error) {
	return []model.DataTopNetworkTraffic{}, nil
}
