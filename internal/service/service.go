package service

import (
	"errors"
	"github.com/filatkinen/sysmon/internal"
	"github.com/filatkinen/sysmon/internal/config"
	"github.com/filatkinen/sysmon/internal/model"
)

type DataElement struct {
	loadAvg           model.DataLoadAvg
	cpuAvgStats       model.DataCpuAvgStats
	disksLoad         []model.DataDisksLoad
	disksUsage        []model.DataDisksUsage
	networkListen     []model.DataNetworkListen
	networkStates     []model.DataNetworkStates
	topNetworkProto   []model.DataTopNetworkProto
	topNetworkTraffic []model.DataTopNetworkTraffic
}

type Data struct {
	element    map[int]*DataElement
	index      []int
	indexCount int
}

type Service struct {
	stat        internal.StatGetter
	conf        config.ServiceConfig
	data        Data
	maxElements int
}

func NewService(serviceConfig config.ServiceConfig, statSource internal.StatGetter) (*Service, error) {
	if serviceConfig.Depth < serviceConfig.ScrapeInterval {
		return nil, errors.New("depth interval cannot be less then scrap interval")
	}
	maxElements := int(serviceConfig.Depth / serviceConfig.ScrapeInterval % 1_000_000_000)

	return &Service{
		stat: statSource,
		conf: serviceConfig,
		data: Data{
			element:    make(map[int]*DataElement, maxElements),
			index:      make([]int, 0, maxElements),
			indexCount: 0,
		},
		maxElements: maxElements,
	}, nil
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) getStatData() error {
	return nil
}

func (s *Service) cleanOldData() error {
	return nil
}
