package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const (
	DefaultPort = "8080"
	DefaultHost = "127.0.0.1"
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

	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type Route struct {
	Path       string         `yaml:"path" json:"path"`
	Method     string         `yaml:"method" json:"method"`
	StatusCode int            `yaml:"statusCode" json:"statusCode"`
	Response   string         `yaml:"response" json:"response"`
	Variants   []RouteVariant `yaml:"variants,omitempty" json:"variants"`
}

type RouteVariant struct {
	Params     map[string]string `yaml:"params" json:"params"`
	StatusCode int               `yaml:"statusCode" json:"statusCode"`
	Response   string            `yaml:"response" json:"response"`
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
