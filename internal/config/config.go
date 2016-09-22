package config

import (
	"log"

	"github.com/ContaAzul/env"
)

// Config config struct
type Config struct {
	User     string   `env:"LIBRATO_EMAIL"`
	Token    string   `env:"LIBRATO_TOKEN"`
	URL      string   `env:"HYSTRIX_URL"`
	Clusters []string `env:"HYSTRIX_CLUSTERS"`
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
