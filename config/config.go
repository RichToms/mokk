package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Name   string  `yaml:"name"`
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Path       string `yaml:"path"`
	Method     string `yaml:"method"`
	StatusCode int    `yaml:"statusCode"`
	Response   string `yaml:"response"`
}

func LoadConfig(cfgPath string) (Config, error) {
	cfgString, err := os.ReadFile(cfgPath)
	if err != nil {
		panic(err)
	}
	var cfg Config

	err = yaml.Unmarshal([]byte(cfgString), &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
