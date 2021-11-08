package config

import (
	"io/ioutil"
	"time"

	"github.com/qiangxue/go-env"
	"gopkg.in/yaml.v2"
	"users-balance-microservice/pkg/log"
)

const defaultServerPort = 8080

// Config represents an application configuration.
type Config struct {
	// the server port. Defaults to 8080.
	ServerPort int `yaml:"server_port" env:"SERVER_PORT"`
	// the expiration time of currency rates. Defaults to 10 minutes.
	RatesExpiration time.Duration `yaml:"rates_expiration"`
	// the data source name (DSN) for connecting to the database. Required.
	DSN string `yaml:"dsn"`
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(file string, logger log.Logger) (*Config, error) {
	// default config
	c := Config{
		ServerPort:      defaultServerPort,
		RatesExpiration: 10 * time.Minute,
	}

	// load from YAML config file
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	// load from environment variables prefixed with "APP_"
	if err = env.New("APP_", logger.Infof).Load(&c); err != nil {
		return nil, err
	}

	return &c, err
}
