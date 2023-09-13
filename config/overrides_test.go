package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolvePortDefaultsToCommandPort(t *testing.T) {
	cfg := Config{}

	port := resolvePort(cfg, DefaultPort)
	assert.Equal(t, DefaultPort, port)
}

func TestResolvePortUsesTheConfigWhenCommandIsDefault(t *testing.T) {
	cfg := Config{
		Options: Options{
			Port: "1234",
		},
	}

	port := resolvePort(cfg, DefaultPort)
	assert.Equal(t, "1234", port)
}

func TestResolvePortUsesTheCommandValueWhenCommandIsNotDefault(t *testing.T) {
	cfg := Config{
		Options: Options{
			Port: "1234",
		},
	}

	port := resolvePort(cfg, "5678")
	assert.Equal(t, "5678", port)
}

func TestResolveHostUsesTheDefaultWhenNoOverridesAreSet(t *testing.T) {
	cfg := Config{}

	host := resolveHost(cfg, "")
	assert.Equal(t, DefaultHost, host)
}

func TestResolveHostUsesTheEnvironmentWhenSet(t *testing.T) {
	cfg := Config{}

	host := resolveHost(cfg, "0.0.0.0")
	assert.Equal(t, "0.0.0.0", host)
}

func TestResolveHostUsesTheConfigWhenNoEnvironmentIsSet(t *testing.T) {
	cfg := Config{
		Options: Options{
			Host: "0.0.0.0",
		},
	}

	host := resolveHost(cfg, "")
	assert.Equal(t, "0.0.0.0", host)
}
