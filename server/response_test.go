package server

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/richtoms/mokk/config"
	"github.com/stretchr/testify/assert"
)

func TestGetResponse(t *testing.T) {
	tests := []struct {
		name          string
		params        map[string]string
		route         config.Route
		assertionFunc func(t *testing.T, res Response)
	}{
		{
			name:   "a route with no params is matched",
			params: map[string]string{},
			route: config.Route{
				Path:       "test",
				Method:     fiber.MethodGet,
				StatusCode: fiber.StatusOK,
				Response:   `{}`,
			},
			assertionFunc: func(t *testing.T, res Response) {
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
				assert.Equal(t, `{}`, res.Response)
			},
		},
		{
			name: "a route with params uses the default when no variants are provided",
			params: map[string]string{
				"user": "123",
			},
			route: config.Route{
				Path:       "users/:user",
				Method:     fiber.MethodGet,
				StatusCode: fiber.StatusOK,
				Response:   `{}`,
			},
			assertionFunc: func(t *testing.T, res Response) {
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
				assert.Equal(t, `{}`, res.Response)
			},
		},
		{
			name: "a route with params uses the found variant",
			params: map[string]string{
				"user": "abc-123",
			},
			route: config.Route{
				Path:       "users/:user",
				Method:     fiber.MethodGet,
				StatusCode: fiber.StatusOK,
				Response:   `{}`,
				Variants: []config.RouteVariant{
					// Test should match the below struct
					{
						Params: map[string]string{
							"user": "abc-123",
						},
						StatusCode: fiber.StatusNotFound,
						Response:   `{"status":"failure"}`,
					},
				},
			},
			assertionFunc: func(t *testing.T, res Response) {
				assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
				assert.Equal(t, `{"status":"failure"}`, res.Response)
			},
		},
		{
			name: "a route with variants returns the default if no variant is matched",
			params: map[string]string{
				"user": "1",
			},
			route: config.Route{
				Path:       "users/:user",
				Method:     fiber.MethodGet,
				StatusCode: fiber.StatusOK,
				Response:   `{}`,
				Variants: []config.RouteVariant{
					{
						Params: map[string]string{
							"user": "abc-123",
						},
						StatusCode: fiber.StatusNotFound,
						Response:   `{"status":"failure"}`,
					},
				},
			},
			assertionFunc: func(t *testing.T, res Response) {
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
				assert.Equal(t, `{}`, res.Response)
			},
		},
		{
			name: "a route with multiple variants returns the correct variant matched",
			params: map[string]string{
				"user": "abcdef",
			},
			route: config.Route{
				Path:       "users/:user",
				Method:     fiber.MethodGet,
				StatusCode: fiber.StatusOK,
				Response:   `{}`,
				Variants: []config.RouteVariant{
					{
						Params: map[string]string{
							"user": "abc-123",
						},
						StatusCode: fiber.StatusNotFound,
						Response:   `{"status":"failure"}`,
					},
					// Test should match the below struct
					{
						Params: map[string]string{
							"user": "abcdef",
						},
						StatusCode: fiber.StatusUnauthorized,
						Response:   `{"status":"auth-failed"}`,
					},
				},
			},
			assertionFunc: func(t *testing.T, res Response) {
				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
				assert.Equal(t, `{"status":"auth-failed"}`, res.Response)
			},
		},
		{
			name: "a route with multiple params only matches the variant when all params match",
			params: map[string]string{
				"user":   "1234",
				"client": "5678",
			},
			route: config.Route{
				Path:       "users/:user/clients/:client",
				Method:     fiber.MethodGet,
				StatusCode: fiber.StatusNotFound,
				Response:   `{}`,
				Variants: []config.RouteVariant{
					{
						Params: map[string]string{
							"user":   "1234",
							"client": "0987",
						},
						StatusCode: fiber.StatusUnauthorized,
						Response:   `{"status":"failure"}`,
					},
					// Test should match the below struct
					{
						Params: map[string]string{
							"user":   "1234",
							"client": "5678",
						},
						StatusCode: fiber.StatusOK,
						Response:   `{"status":"success"}`,
					},
				},
			},
			assertionFunc: func(t *testing.T, res Response) {
				assert.Equal(t, fiber.StatusOK, res.StatusCode)
				assert.Equal(t, `{"status":"success"}`, res.Response)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := getResponse(test.params, test.route)

			test.assertionFunc(t, res)
		})
	}
}
