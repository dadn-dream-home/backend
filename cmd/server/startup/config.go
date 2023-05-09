package startup

import (
	"context"

	"github.com/dadn-dream-home/x/server/telemetry"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	ServerConfig   `hcl:"server,block"`
	DatabaseConfig `hcl:"database,block"`
	MQTTConfig     `hcl:"mqtt,block"`
}

type ServerConfig struct {
	Port int `hcl:"port"`
}

type DatabaseConfig struct {
	ConnectionString string `hcl:"connection_string"`
	MigrationsPath   string `hcl:"migrations_path"`
}

type MQTTConfig struct {
	Brokers []string `hcl:"brokers"`
}

func OpenConfig(ctx context.Context) *Config {
	log := telemetry.GetLogger(ctx)

	path := "config.hcl"
	log = log.WithField("path", path)

	var config Config
	err := hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		log.WithError(err).Fatal("failed to load config")
	}

	log.Info("loaded config")

	return &config
}
