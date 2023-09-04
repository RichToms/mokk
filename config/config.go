package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Name    string  `yaml:"name"`
	Options Options `yaml:"options"`
	Routes  []Route `yaml:"routes"`
}

type Options struct {
	// When a request body is received, attempt to print the value to the console
	// after the table record.
	PrintRequestBody bool `yaml:"printRequestBody"`
}

type Route struct {
	Path       string         `yaml:"path"`
	Method     string         `yaml:"method"`
	StatusCode int            `yaml:"statusCode"`
	Response   string         `yaml:"response"`
	Variants   []RouteVariant `yaml:"variants,omitempty"`
}

type RouteVariant struct {
	Params     map[string]string `yaml:"params"`
	StatusCode int               `yaml:"statusCode"`
	Response   string            `yaml:"response"`
}

// LoadConfigFromFile attempts to load the config from the given file path.
func LoadConfigFromFile(cfgPath string) (Config, error) {
	cfgString, err := os.ReadFile(cfgPath)
	if err != nil {
		panic(err)
	}

	return loadConfigString(string(cfgString))
}

// loadConfigString attempts to unmarshal a YAML string into a Mokk config.
func loadConfigString(cfgString string) (Config, error) {
	var cfg Config

	err := yaml.Unmarshal([]byte(cfgString), &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
