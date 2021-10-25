package config

import (
	"io/ioutil"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"users-balance-microservice/pkg/log"
)

const defaultServerPort = 8080

// Config represents an application configuration.
type Config struct {
	// the server port. Defaults to 8080
	ServerPort int `yaml:"server_port" env:"SERVER_PORT"`
	// the database url.
	StorageUrl string `yaml:"storage_url"`
	// the storage driver name
	Driver string `yaml:"driver"`
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(file string, logger log.Logger) (*Config, error) {
	// default config
	c := Config{
		ServerPort: defaultServerPort,
		StorageUrl: "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		Driver:     "pgx",
	}

	// load from YAML config file
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	/* load from environment variables prefixed with "APP_"
	if err = env.New("APP_", logger.Infof).Load(&c); err != nil {
		return nil, err
	}*/

	return &c, err
}
