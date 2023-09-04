package server

import "github.com/richtoms/mokk/config"

type Response struct {
	StatusCode int
	Response   string
}

// getResponse attempts to find the correct response based on the request.
func getResponse(params map[string]string, route config.Route) Response {
	res := Response{
		StatusCode: route.StatusCode,
		Response:   route.Response,
	}

	if len(route.Variants) > 0 {
		for _, variant := range route.Variants {
			matches := make([]bool, 0)
			for key, value := range variant.Params {
				if params[key] == value {
					matches = append(matches, true)
				}
			}

			if len(matches) == len(variant.Params) {
				res = Response{
					StatusCode: variant.StatusCode,
					Response:   variant.Response,
				}
			}
		}
	}

	return res
}
