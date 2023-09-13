package server

import "fmt"

type Options struct {
	Host string
	Port string
}

func (o Options) resolveHost() string {
	return fmt.Sprintf("%s:%s", o.Host, o.Port)
}
