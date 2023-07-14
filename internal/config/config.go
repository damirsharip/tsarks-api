package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	configYml = "config.yml"
)

type Gateway struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type PgxPool struct {
	// MaxConnLifetime is the duration since creation after which a connection will be automatically closed.
	MaxConnLifetime int64 `yaml:"max_conn_lifetime"`
	// MaxConnIdleTime is the duration after which an idle connection will be automatically closed by the health check.
	MaxConnIdleTime int64 `yaml:"max_conn_idle_time"`
	// MaxConns is the maximum size of the pool. The default is the greater of 4 or runtime.NumCPU().
	MaxConns int32 `yaml:"max_conns"`
	// MinConns is the minimum size of the pool. The health check will increase the number of connections to this
	// amount if it had dropped below.
	MinConns int32 `yaml:"min_conns"`
}

type Database struct {
	Dsn string `yaml:"dsn"`
}

type Config struct {
	PgxPool  PgxPool  `yaml:"pgxpool"`
	Gateway  Gateway  `yaml:"gateway"`
	Database Database `yaml:"database"`
}

var (
	cfg *Config
)

func ReadConfigYaml() (err error) {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(configYml)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return err
	}

	return
}

func Get() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}
