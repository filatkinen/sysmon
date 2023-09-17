package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	DefaultScrapeInterval = time.Second * 5
	DefaultCleanInterval  = time.Minute * 5
	DefaultDepth          = time.Hour * 1
	DefaultPort           = "50051"
	DefaultAddress        = "0.0.0.0"
)

type ServiceConfig struct {
	LA          bool
	AvgCPU      bool
	DisksLoad   bool
	DisksUse    bool
	NetworkTop  bool
	NetworkStat bool

	ScrapeInterval time.Duration
	CleanInterval  time.Duration
	Depth          time.Duration

	Port    string
	Address string
}

func NewConfig(in string) (ServiceConfig, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(in)
	if err := viper.ReadInConfig(); err != nil {
		return ServiceConfig{}, err
	}

	viper.SetDefault("subsystems.la", true)
	viper.SetDefault("subsystems.avgcpu", true)
	viper.SetDefault("subsystems.disksload", true)
	viper.SetDefault("subsystems.disksinfo", true)
	viper.SetDefault("subsystems.networktop", true)
	viper.SetDefault("subsystems.networkstat", true)

	viper.SetDefault("scrapeinterval", DefaultScrapeInterval.String())

	viper.SetDefault("bindings.port", DefaultPort)
	viper.SetDefault("bindings.address", DefaultAddress)

	duration, err := time.ParseDuration(viper.GetString("scrapeinterval"))
	if err != nil {
		log.Printf("Error parsing scrapeinterval value: %s, using defaul tvalue:%s", err.Error(), DefaultScrapeInterval)
		duration = DefaultScrapeInterval
	}

	depth, err := time.ParseDuration(viper.GetString("depth"))
	if err != nil {
		log.Printf("Error parsing scrapeinterval value: %s, using defaul tvalue:%s", err.Error(), DefaultDepth)
		depth = DefaultDepth
	}

	clean, err := time.ParseDuration(viper.GetString("cleaninterval"))
	if err != nil {
		log.Printf("Error parsing cleaninterval value: %s, using defaul tvalue:%s", err.Error(), DefaultCleanInterval)
		clean = DefaultCleanInterval
	}

	config := ServiceConfig{
		LA:             viper.GetBool("subsystems.la"),
		AvgCPU:         viper.GetBool("subsystems.avgcpu"),
		DisksLoad:      viper.GetBool("subsystems.disksload"),
		DisksUse:       viper.GetBool("subsystems.disksuse"),
		NetworkTop:     viper.GetBool("subsystems.networktop"),
		NetworkStat:    viper.GetBool("subsystems.networkstat"),
		Port:           viper.GetString("bindings.port"),
		Address:        viper.GetString("bindings.address"),
		ScrapeInterval: duration,
		Depth:          depth,
		CleanInterval:  clean,
	}

	return config, nil
}
