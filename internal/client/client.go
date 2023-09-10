package client

import (
	"context"
	"encoding/json"
	"errors"
	pb "github.com/filatkinen/sysmon/internal/grpc/sysmon"
	"github.com/filatkinen/sysmon/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

type DataSlice struct {
	slice []Element
}

type Element struct {
	header string
	data   [][]string
}

type ClientSysMon struct {
	client    pb.SysmonDataClient
	conn      *grpc.ClientConn
	addr      string
	every_m   int
	average_n int
}

func NewClientSysMon(address string, every_m, average_n int) *ClientSysMon {
	return &ClientSysMon{
		addr:      address,
		every_m:   every_m,
		average_n: average_n,
	}
}

func (s *ClientSysMon) Start() error {
	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	s.conn = conn
	s.client = pb.NewSysmonDataClient(conn)
	return nil
}

func (s *ClientSysMon) Close() error {
	return s.conn.Close()
}

func (s *ClientSysMon) GetData(f func(data []model.DataToClientStamp)) error {
	ctx := context.Background()
	DataStream, err := s.client.SendSysmonDataToClient(ctx, &pb.QueryParam{
		EveryM:   int32(s.every_m),
		AverageN: int32(s.average_n),
	})
	if err != nil {
		return err
	}
	for {
		dataStream, err := DataStream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		var m []model.DataToClientStamp
		err = json.Unmarshal(dataStream.Data, &m)
		if err != nil {
			return err
		}
		f(m)
	}
	return nil
}
