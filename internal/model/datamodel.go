package model

import (
	"sync"
)

type NameHeader struct {
	Name            string
	Header          []string
	SortHeaderField int
	SortDescending  bool
}

var StampNameHeaders = []NameHeader{
	{
		Name:            "LoadAvg",
		Header:          append([]string(nil), "5 min", "10 min", "15 min"),
		SortHeaderField: 0,
		SortDescending:  false,
	},
	{
		Name:            "CPUAvgStats",
		Header:          append([]string(nil), "User", "System", "Idle"),
		SortHeaderField: 0,
		SortDescending:  false,
	},
	{
		Name:            "DisksLoad",
		Header:          append([]string(nil), "Device", "Tps", "Kbps"),
		SortHeaderField: 1,
		SortDescending:  true,
	},
	{
		Name: "DisksUsage",
		Header: append([]string(nil), "Mount Point", "File System",
			"Usage Inodes", "Usage Inode Percent", "Usage Mb", "Usage Mb Percent"),
		SortHeaderField: 0,
		SortDescending:  false,
	},
	{
		Name:            "NetworkListen",
		Header:          append([]string(nil), "Protocol", "Command", "PID", "USER", "PORT"),
		SortHeaderField: 0,
		SortDescending:  false,
	},
	{
		Name:            "NetworkStates",
		Header:          append([]string(nil), "Status", "Number"),
		SortHeaderField: 0,
		SortDescending:  false,
	},
	{
		Name:            "TopNetworkProto",
		Header:          append([]string(nil), "Protocol", "Bytes", "Bytes Percent"),
		SortHeaderField: 2,
		SortDescending:  true,
	},
	{
		Name:            "TopNetworkTraffic",
		Header:          append([]string(nil), "Source", "Destination", "Protocol", "Bps"),
		SortHeaderField: 3,
		SortDescending:  true,
	},
}

type Element struct {
	CountAble    bool
	StringField  string
	NumberField  float64
	DecimalField int
	Count        int
}

type ElMapType map[string][]Element

type StampsElements struct {
	ElMap               ElMapType
	IdxStampNameHeaders int
}

type StampsData struct {
	Data []StampsElements
}

type Data struct {
	Elements    map[int]StampsData
	Index       []int
	Counter     int
	MaxElements int
	Lock        sync.RWMutex
}

type DataToClientStamp struct {
	Name      string
	IdxHeader int
	Data      [][]string
}
