#!/bin/bash

#apt install  protobuf-compiler
#go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

protoc -I  ./internal/grpc/ ./internal/grpc/sysmon.proto   --go_out=./internal/grpc/
protoc -I  ./internal/grpc/ ./internal/grpc/sysmon.proto   --go-grpc_out=require_unimplemented_servers=false:./internal/grpc/


#package example

//
//import (
//	"fmt"
//	"github.com/rafacas/sysstats"
//)
//
//func main() {
//	s, err := sysstats.GetLoadAvg()
//	if err != nil {
//		return
//	}
//	fmt.Println(s)
//
//	stats, err := sysstats.GetCpuStatsInterval(1)
//	fmt.Println(stats)
//}
