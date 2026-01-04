package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type Config struct {
	Server   Server
	Log      Log
	Database Database
}

type Server struct {
	ListenAddress *string `yaml:"listen_address" validate:"omitempty,hostname_port"`
}

type Log struct {
	Level  *string `validate:"omitempty,oneofci=trace debug info warn error fatal panic"`
	Format *string `validate:"omitempty,oneof=console json"`
}

type Database struct {
	Host     *string `validate:"omitnil,hostname|ip"`
	Port     *uint16 `validate:"omitnil,port"`
	Database *string `validate:"omitempty,min=1"`
	User     *string `validate:"omitempty,min=1"`
	Password *string `validate:"omitempty,min=1"`
}

func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config file reading: %w", err)
	}

	dec := yaml.NewDecoder(
		bytes.NewReader(raw),
		yaml.Strict(),
		yaml.DisallowUnknownField(),
	)

	var ret Config
	setDefaults(&ret)

	if err := dec.Decode(&ret); err != nil {
		return nil, fmt.Errorf("config file parsing: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&ret); err != nil {
		return nil, fmt.Errorf("config file parsing: %w", err)
	}

	return &ret, nil
}

func setDefaults(cfg *Config) {
	if cfg.Server.ListenAddress == nil {
		defaultAddr := ":8080"
		cfg.Server.ListenAddress = &defaultAddr
	}
	if cfg.Log.Level == nil {
		defaultLevel := "info"
		cfg.Log.Level = &defaultLevel
	}
	if cfg.Log.Format == nil {
		defaultFormat := "json"
		cfg.Log.Format = &defaultFormat
	}
}
