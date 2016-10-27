package config

import (
	"log"

	"github.com/caarlos0/env"
)

// Config config struct
type Config struct {
	User            string   `env:"LIBRATO_EMAIL,required"`
	Token           string   `env:"LIBRATO_TOKEN,required"`
	URL             string   `env:"HYSTRIX_URL,required"`
	Clusters        []string `env:"HYSTRIX_CLUSTERS"`
	ReportLatencies []string `env:"HYSTRIX_REPORT_LATENCIES" envSeparator:"," envDefault:"100th,99.5th,99th,95th,90th,75th,50th,25th,0th,mean"`
	ReportInterval  float64  `env:"HYSTRIX_REPORT_INTERVAL" envDefault:"5"`
}

// Get the config
func Get() Config {
	var conf Config
	err := env.Parse(&conf)
	if err != nil {
		log.Fatalln(err)
	}
	return conf
}
