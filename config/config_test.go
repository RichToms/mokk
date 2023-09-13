package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	cfg, err := LoadConfigFromFile("../mokk.yml")

	assert.NoError(t, err)
	assert.Equal(t, "Mokk Users Server", cfg.Name)
}

func TestLoadConfigString(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		assertionFunc func(t *testing.T, cfg Config)
	}{
		{
			name:  "server name can be set",
			input: nameIsReadInput,
			assertionFunc: func(t *testing.T, cfg Config) {
				assert.Equal(t, "Mokk Server", cfg.Name)
			},
		},
		{
			name:  "options can be set",
			input: configWithOptions,
			assertionFunc: func(t *testing.T, cfg Config) {
				assert.Equal(t, true, cfg.Options.PrintRequestBody)
				assert.Equal(t, "9001", cfg.Options.Port)
				assert.Equal(t, "0.0.0.0", cfg.Options.Host)
			},
		},
		{
			name:  "routes can be set (with variants)",
			input: configWithRouteVariants,
			assertionFunc: func(t *testing.T, cfg Config) {
				assert.Len(t, cfg.Routes, 2)

				var r Route
				var v RouteVariant

				r = cfg.Routes[0]
				assert.Equal(t, "users/:user", r.Path)
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, 200, r.StatusCode)
				assert.Equal(t, "{\"status\":\"success\"}", r.Response)

				assert.Len(t, r.Variants, 1)
				v = r.Variants[0]

				assert.Equal(t, "123", v.Params["user"])
				assert.Equal(t, 404, v.StatusCode)
				assert.Equal(t, "{\"status\":\"failure\"}", v.Response)

				r = cfg.Routes[1]
				assert.Equal(t, "users/:user/clients/:client", r.Path)
				assert.Equal(t, "PATCH", r.Method)
				assert.Equal(t, 200, r.StatusCode)
				assert.Equal(t, "{\"status\":\"success\"}", r.Response)

				assert.Len(t, r.Variants, 2)
				v = r.Variants[0]
				assert.Equal(t, "123", v.Params["user"])
				assert.Equal(t, "123", v.Params["client"])
				assert.Equal(t, 200, v.StatusCode)
				assert.Equal(t, "{\"status\":\"success\",\"client\":{}}", v.Response)

				v = r.Variants[1]
				assert.Equal(t, "123", v.Params["user"])
				assert.Equal(t, "456", v.Params["client"])
				assert.Equal(t, 404, v.StatusCode)
				assert.Equal(t, "{\"status\":\"failure\"}", v.Response)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := loadConfigString(test.input)
			assert.NoError(t, err)

			test.assertionFunc(t, cfg)
		})
	}
}

const nameIsReadInput = `
name: Mokk Server
`

const configWithOptions = `
name: Mokk Server with options

options:
  printRequestBody: true
  port: 9001
  host: 0.0.0.0

routes:
  - path: example
    method: GET
    statusCode: 200
    response: '{}'
`

const configWithRouteVariants = `
name: Mokk Server with routes (variants)

routes:
  - path: users/:user
    method: GET
    statusCode: 200
    response: '{"status":"success"}'
    variants:
      - params:
          user: 123
        statusCode: 404
        response: '{"status":"failure"}'
  - path: users/:user/clients/:client
    method: PATCH
    statusCode: 200
    response: '{"status":"success"}'
    variants:
      - params:
          user: 123
          client: 123
        statusCode: 200
        response: '{"status":"success","client":{}}'
      - params:
          user: 123
          client: 456
        statusCode: 404
        response: '{"status":"failure"}'
`
