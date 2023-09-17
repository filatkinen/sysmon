package main

import (
	"fmt"
	"github.com/rafacas/sysstats"
)

func main() {
	s, err := sysstats.GetLoadAvg()
	if err != nil {
		return
	}
	fmt.Println(s)

	stat1, err := sysstats.GetDiskUsage()
	fmt.Println(stat1)

	//stats, err := sysstats.GetCpuStatsInterval(1)
	//fmt.Println(stats)
}
