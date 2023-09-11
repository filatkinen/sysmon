package service

import (
	"errors"
	"github.com/filatkinen/sysmon/internal"
	"github.com/filatkinen/sysmon/internal/config"
	pb "github.com/filatkinen/sysmon/internal/grpc/sysmon"
	"github.com/filatkinen/sysmon/internal/model"
	"google.golang.org/grpc"
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

type Service struct {
	stat     internal.StatGetter
	conf     config.ServiceConfig
	data     model.Data
	conn     *grpc.Server
	exitChan chan struct{}
	wg       sync.WaitGroup
}

func NewService(serviceConfig config.ServiceConfig, statSource internal.StatGetter) (*Service, error) {
	if serviceConfig.Depth < serviceConfig.ScrapeInterval {
		return nil, errors.New("depth interval cannot be less then scrap interval")
	}
	maxElements := int((serviceConfig.Depth / serviceConfig.ScrapeInterval)) + 1
	log.Printf("Creating sysmon service: %+v\n", serviceConfig)
	return &Service{
		stat:     statSource,
		conf:     serviceConfig,
		exitChan: make(chan struct{}),
		wg:       sync.WaitGroup{},
		data: model.Data{
			Elements:    make(map[int]model.StampsData, maxElements),
			Index:       make([]int, 0, maxElements),
			MaxElements: maxElements,
		},
	}, nil
}

func (s *Service) Start() error {
	log.Printf("Starting sysmon service...\n")
	wg := sync.WaitGroup{}
	wg.Add(3)
	var err error
	go func() {
		defer wg.Done()
		log.Printf("Starting sysmon service. GRPC subsystem...\n")
		err = s.startGRPC()
		if err != nil {
			log.Printf("Failed to start GRPC server: %s ", err.Error())
		}
	}()
	go func() {
		defer wg.Done()
		s.getStatData()
	}()
	go func() {
		defer wg.Done()
		s.cleanOldData()
	}()
	wg.Wait()
	return err
}

func (s *Service) startGRPC() error {
	s.conn = grpc.NewServer()
	lis, err := net.Listen("tcp", net.JoinHostPort(s.conf.Address, s.conf.Port))
	if err != nil {
		return err
	}
	pb.RegisterSysmonDataServer(s.conn, s)
	if err := s.conn.Serve(lis); err != nil {
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return err
		}
	}
	return nil
}

func (s *Service) Stop() error {
	log.Printf("Stopping sysmon service...\n")
	log.Printf("Stopping sysmon service. GRPC subsystem...\n")
	close(s.exitChan)
	s.wg.Wait()
	s.conn.Stop()
	return nil
}

func (s *Service) getStatData() {
	ticker := time.NewTicker(s.conf.ScrapeInterval)
	defer func() {
		ticker.Stop()
	}()

	type statFunc = func() (model.ElMapType, error)
	var wgService sync.WaitGroup
	getstat := func(f statFunc) model.ElMapType {
		defer wgService.Done()
		m, err := f()
		if err != nil {
			funcname := "unknown"
			pc, _, _, ok := runtime.Caller(0)
			if ok {
				me := runtime.FuncForPC(pc)
				if me != nil {
					funcname = me.Name()
				}
			}
			log.Printf("got error quering using function %s, error: %s", funcname, err)
			return nil
		}
		return m
	}

	addstat := func(sd *model.StampsData, se model.ElMapType, idx int) {
		if se != nil {
			var sedata model.StampsElements
			sedata.ElMap = se
			sedata.IdxStampNameHeaders = idx
			sd.Data = append(sd.Data, sedata)
		}
	}

	collectStat := func() {
		var LoadAvg, CpuAvgStats, DisksLoad, DisksInfo, NetworkListen, NetworkStates, TopNetworkProto, TopNetworkTraffic model.ElMapType
		switch {
		case s.conf.LA:
			wgService.Add(1)
			go func() { LoadAvg = getstat(s.stat.LoadAvg) }()
			fallthrough
		case s.conf.AvgCpu:
			wgService.Add(1)
			go func() { CpuAvgStats = getstat(s.stat.CpuAvgStats) }()
			fallthrough
		case s.conf.DisksInfo:
			wgService.Add(1)
			go func() { DisksInfo = getstat(s.stat.DisksUsage) }()
			fallthrough
		case s.conf.DisksLoad:
			wgService.Add(1)
			go func() { DisksLoad = getstat(s.stat.DisksLoad) }()
			fallthrough
		case s.conf.NetworkStat:
			wgService.Add(2)
			go func() { NetworkListen = getstat(s.stat.NetworkListen) }()
			go func() { NetworkStates = getstat(s.stat.NetworkStates) }()
			fallthrough
		case s.conf.NetworkTop:
			wgService.Add(2)
			go func() { TopNetworkProto = getstat(s.stat.TopNetworkProto) }()
			go func() { TopNetworkTraffic = getstat(s.stat.TopNetworkTraffic) }()
		}

		wgService.Wait()
		var sd model.StampsData

		switch {
		case s.conf.LA:
			addstat(&sd, LoadAvg, 0)
			fallthrough
		case s.conf.AvgCpu:
			addstat(&sd, CpuAvgStats, 1)
			fallthrough
		case s.conf.DisksInfo:
			addstat(&sd, DisksInfo, 2)
			fallthrough
		case s.conf.DisksLoad:
			addstat(&sd, DisksLoad, 3)
			fallthrough
		case s.conf.NetworkStat:
			addstat(&sd, NetworkListen, 4)
			addstat(&sd, NetworkStates, 5)
			fallthrough
		case s.conf.NetworkTop:
			addstat(&sd, TopNetworkProto, 6)
			addstat(&sd, TopNetworkTraffic, 7)
		}
		s.data.Lock.Lock()
		s.data.Counter++
		s.data.Index = append(s.data.Index, s.data.Counter)
		s.data.Elements[s.data.Counter] = sd
		s.data.Lock.Unlock()
	}

	collectStat()
	for {
		select {
		case <-ticker.C:
			collectStat()
		case <-s.exitChan:
			return
		}
	}
}

func (s *Service) cleanOldData() {
	ticker := time.NewTicker(s.conf.CleanInterval)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			s.data.Lock.Lock()
			num := len(s.data.Index) - s.data.MaxElements
			if num > 0 {
				for i := 0; i < num; i++ {
					idx := s.data.Index[0]
					delete(s.data.Elements, idx)
					s.data.Index = s.data.Index[1:len(s.data.Index)]
				}
			}
			s.data.Lock.Unlock()
		case <-s.exitChan:
			return
		}
	}
}

func (s *Service) countDataClient(averageN int) (*model.StampsData, bool, error) {
	s.data.Lock.RLock()
	defer s.data.Lock.RUnlock()

	quantityElements := len(s.data.Index)
	scrapeSeconds := int(s.conf.ScrapeInterval / time.Second)
	if scrapeSeconds*quantityElements < averageN {
		return nil, false, nil
	}

	quantityCl := averageN / scrapeSeconds
	if quantityCl == 0 {
		quantityCl = 1
	}
	lastdata := s.data.Elements[s.data.Index[quantityElements-1]]
	clientData := copyStampsData(lastdata)

	for count := 2; count <= quantityCl; count++ {
		clientData = sumStampsData(*clientData, s.data.Elements[s.data.Index[quantityElements-count]])
	}
	averageStampsData(clientData, quantityCl)
	return clientData, true, nil
}

func copyStampsData(source model.StampsData) *model.StampsData {
	var m model.StampsData
	m.Data = make([]model.StampsElements, len(source.Data))
	for i := range source.Data {
		m.Data[i].IdxStampNameHeaders = source.Data[i].IdxStampNameHeaders
		m.Data[i].ElMap = make(map[string][]model.Element, len(source.Data[i].ElMap))
		for k := range source.Data[i].ElMap {
			m.Data[i].ElMap[k] = append([]model.Element(nil), source.Data[i].ElMap[k]...)
		}
	}
	return &m
}

func sumStampsData(s1 model.StampsData, s2 model.StampsData) *model.StampsData {
	var m model.StampsData
	m.Data = make([]model.StampsElements, len(s1.Data), len(s1.Data))
	for i := range s1.Data {
		hashI := make(map[string]bool, len(s1.Data[i].ElMap))
		for k := range s1.Data[i].ElMap {
			hashI[k] = true
		}
		for k := range s2.Data[i].ElMap {
			hashI[k] = true
		}

		m.Data[i].IdxStampNameHeaders = s1.Data[i].IdxStampNameHeaders
		m.Data[i].ElMap = make(map[string][]model.Element, len(s1.Data[i].ElMap))
		for k := range hashI {
			_, ok1 := s1.Data[i].ElMap[k]
			_, ok2 := s2.Data[i].ElMap[k]
			switch {
			case ok1 && ok2:
				m.Data[i].ElMap[k] = make([]model.Element, 0, len(s1.Data[i].ElMap[k]))
				for idx, v := range s1.Data[i].ElMap[k] {
					if v.CountAble {
						v.NumberField += s2.Data[i].ElMap[k][idx].NumberField
					}
					m.Data[i].ElMap[k] = append(m.Data[i].ElMap[k], v)
				}
			case ok1:
				m.Data[i].ElMap[k] = append([]model.Element(nil), s1.Data[i].ElMap[k]...)
			case ok2:
				m.Data[i].ElMap[k] = append([]model.Element(nil), s2.Data[i].ElMap[k]...)
			}
		}
	}
	return &m
}

func averageStampsData(s *model.StampsData, count int) {
	for i := range s.Data {
		for k := range s.Data[i].ElMap {
			for idx := range s.Data[i].ElMap[k] {
				if s.Data[i].ElMap[k][idx].CountAble {
					s.Data[i].ElMap[k][idx].NumberField = s.Data[i].ElMap[k][idx].NumberField / float64(count)
				}
			}
		}
	}
}
