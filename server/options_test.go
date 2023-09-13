package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveHost(t *testing.T) {
	tests := []struct {
		host     string
		port     string
		expected string
	}{
		{"127.0.0.1", "8000", "127.0.0.1:8000"},
		{"127.0.0.1", "8080", "127.0.0.1:8080"},
		{"0.0.0.0", "80", "0.0.0.0:80"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s[%s]", test.host, test.port), func(t *testing.T) {
			opt := Options{
				Host: test.host,
				Port: test.port,
			}

			output := opt.resolveHost()
			assert.Equal(t, test.expected, output)
		})
	}
}
