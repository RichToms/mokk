package server

import "fmt"

type Options struct {
	Host string
	Port string
}

// resolveHost builds the host required by the server to understand which host and port to listen on.
func (o Options) resolveHost() string {
	return fmt.Sprintf("%s:%s", o.Host, o.Port)
}
