package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	pb "github.com/filatkinen/sysmon/internal/grpc/sysmon"
	"github.com/filatkinen/sysmon/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SysMon struct {
	client   pb.SysmonDataClient
	conn     *grpc.ClientConn
	addr     string
	everyM   int
	averageN int
}

func NewClient(address string, everyM, averageN int) *SysMon {
	return &SysMon{
		addr:     address,
		everyM:   everyM,
		averageN: averageN,
	}
}

func (s *SysMon) Start() error {
	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	s.conn = conn
	s.client = pb.NewSysmonDataClient(conn)
	return nil
}

func (s *SysMon) Close() error {
	return s.conn.Close()
}

func (s *SysMon) GetData(f func([]model.DataToClientStamp)) error {
	ctx := context.Background()
	DataStream, err := s.client.SendSysmonDataToClient(ctx, &pb.QueryParam{
		EveryM:   int32(s.everyM),
		AverageN: int32(s.averageN),
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
