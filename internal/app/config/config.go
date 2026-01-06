package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

const DefaultConfigPath = "/var/opt/demo-backend-go/config.yaml"

type Config struct {
	Server   Server
	Log      Log
	Database Database
}

type Server struct {
	ListenAddress *string `yaml:"listen_address" validate:"omitempty,hostname_port"`
}

type Log struct {
	Level  string `validate:"omitempty,oneofci=trace debug info warn error fatal panic"`
	Format string `validate:"omitempty,oneof=console json"`
}

type Database struct {
	DriverName     string `yaml:"driver_name" validate:"required"`
	DataSourceName string `yaml:"data_source_name" validate:"required"`
}

func Load() (*Config, error) {
	raw, err := os.ReadFile(DefaultConfigPath)
	if err != nil {
		return nil, fmt.Errorf("config file reading: %v", err)
	}

	dec := yaml.NewDecoder(
		bytes.NewReader(raw),
		yaml.Strict(),
		yaml.DisallowUnknownField(),
	)

	var ret Config
	setDefaults(&ret)

	if err := dec.Decode(&ret); err != nil {
		return nil, fmt.Errorf("config file parsing: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(&ret); err != nil {
		return nil, fmt.Errorf("config file validation: %v", err)
	}

	return &ret, nil
}

func setDefaults(cfg *Config) {
	if cfg.Server.ListenAddress == nil {
		defaultAddr := ":8080"
		cfg.Server.ListenAddress = &defaultAddr
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "json"
	}
}
