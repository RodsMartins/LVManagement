package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port              string `envconfig:"PORT" default:":4000"`
	SessionCookieName string `envconfig:"SESSION_COOKIE_NAME" default:"session"`
	DATABASE_URL      string `envconfig:"DATABASE_URL" default:"postgres://postgres:postgres@172.21.0.3:5432/lvm"`
}

func loadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func MustLoadConfig() *Config {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
