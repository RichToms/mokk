package server

import (
	"slices"
	"strings"

	"github.com/richtoms/mokk/config"
)

type Response struct {
	StatusCode int    `json:"statusCode"`
	Response   string `json:"response"`
}

// getResponse attempts to find the correct response based on the request.
func getResponse(params map[string]string, headers map[string][]string, route config.Route) Response {
	res := Response{
		StatusCode: route.StatusCode,
		Response:   route.Response,
	}

	if len(route.Variants) > 0 {
		for _, variant := range route.Variants {
			paramMatches := make([]bool, 0)
			for key, value := range variant.Params {
				if params[key] == value {
					paramMatches = append(paramMatches, true)
				}
			}

			headerMatches := make([]bool, 0)
			for key, values := range variant.RequestHeaders {
				lowerKey := strings.ToLower(key)
				for _, value := range values {
					if slices.Contains(headers[lowerKey], value) {
						headerMatches = append(headerMatches, true)
					}
				}
			}

			if len(paramMatches) == len(variant.Params) && len(headerMatches) == variant.RequestHeaders.Count() {
				res = Response{
					StatusCode: variant.StatusCode,
					Response:   variant.Response,
				}
			}
		}
	}

	return res
}
