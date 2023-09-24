package service

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"time"

	pb "github.com/filatkinen/sysmon/internal/grpc/sysmon"
	"github.com/filatkinen/sysmon/internal/model"
	"google.golang.org/grpc/peer"
)

func (s *Service) SendSysmonDataToClient(param *pb.QueryParam,
	server pb.SysmonData_SendSysmonDataToClientServer,
) error {
	ticker := time.NewTicker(time.Second * time.Duration(param.EveryM))
	s.wg.Add(1)
	defer func() {
		ticker.Stop()
		s.wg.Done()
	}()

	r, ok := peer.FromContext(server.Context())
	address := ""
	if ok {
		address = r.Addr.String()
	}
	log.Printf("new client GRPC:%s. Query params: average %d seconds, query every %d seconds \n",
		address, param.AverageN, param.EveryM)
	defer log.Printf("end GRPC connection from client:%s\n", address)

	sendData := func() error {
		var dataPB pb.Data
		dataSevice, dataReady := s.CountDataClient(int(param.AverageN))
		if dataReady {
			r := prepareData(dataSevice)
			b, e := json.Marshal(&r)
			if e != nil {
				return e
			}
			dataPB.Data = b
			return server.Send(&dataPB)
		}
		return nil
	}
	err := sendData()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			err = sendData()
			if err != nil {
				return err
			}
		case <-s.exitChan:
			return nil
		}
	}
}

func prepareData(m *model.StampsData) []model.DataToClientStamp {
	result := make([]model.DataToClientStamp, 0, len(m.Data))
	for i := range m.Data {
		data := prepareDataStamp(m.Data[i].ElMap, m.Data[i].IdxStampNameHeaders)
		result = append(result, data)
	}
	return result
}

func prepareDataStamp(data model.ElMapType, headerIdx int) model.DataToClientStamp {
	var result model.DataToClientStamp
	m := make([][]string, 0, len(data)+1)
	Headers := model.StampNameHeaders[headerIdx]
	SortDescending := Headers.SortDescending
	sortIdx := Headers.SortHeaderField
	type sortEl struct {
		d   model.Element
		key string
	}
	sortslice := make([]sortEl, 0, len(data))
	for k := range data {
		sortslice = append(sortslice, struct {
			d   model.Element
			key string
		}{d: data[k][sortIdx], key: k})
	}

	sort.Slice(sortslice, func(i, j int) bool {
		if sortslice[i].d.CountAble {
			if SortDescending {
				return sortslice[i].d.NumberField > sortslice[j].d.NumberField
			}
			return sortslice[i].d.NumberField < sortslice[j].d.NumberField
		}
		if SortDescending {
			return sortslice[i].d.StringField > sortslice[j].d.StringField
		}
		return sortslice[i].d.StringField < sortslice[j].d.StringField
	})

	l := append([]string(nil), Headers.Header...)
	m = append(m, l)
	for k := range sortslice {
		key := sortslice[k].key
		line := make([]string, 0, len(data[key]))
		for _, v := range data[key] {
			if v.CountAble {
				line = append(line, floatToString(v.NumberField, v.DecimalField))
				continue
			}
			line = append(line, v.StringField)
		}
		m = append(m, line)
	}
	result.Data = m
	result.Name = Headers.Name
	result.IdxHeader = headerIdx
	return result
}

// default return with 2 digit after point.
func floatToString(number float64, precision ...int) string {
	if len(precision) > 0 {
		return strconv.FormatFloat(number, 'f', precision[0], 64)
	}
	return strconv.FormatFloat(number, 'f', 2, 64)
}
