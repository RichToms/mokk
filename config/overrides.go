package config

import (
	"github.com/spf13/cobra"
	"os"
)

func (c *Config) OverrideFromCommand(cmd *cobra.Command) {
	c.Options.Port = resolvePort(*c, cmd)
	c.Options.Host = resolveHost(*c, cmd)
}

// resolvePort attempts to determine which port to host the server on.
func resolvePort(cfg Config, cmd *cobra.Command) string {
	cmdPort := cmd.Flag("port").Value.String()

	// If the config has a set port and the command is providing the default port,
	// use the port from the config file
	if cfg.Options.Port != "" && cmdPort == DefaultPort {
		return cfg.Options.Port
	}

	return cmdPort
}

// resolveHost attempts to determine which host to serve on.
func resolveHost(cfg Config, cmd *cobra.Command) string {
	if envHost := os.Getenv("SERVER_HOST"); envHost != "" {
		return envHost
	}

	if cfg.Options.Host != "" {
		return cfg.Options.Host
	}

	return DefaultHost
}
