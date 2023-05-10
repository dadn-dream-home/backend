package startup

import (
	"context"

	"github.com/dadn-dream-home/x/server/telemetry"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"go.uber.org/zap"
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
	log = log.With(zap.String("path", path))

	var config Config
	err := hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	log.Info("loaded config")

	return &config
}
