package server

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/richtoms/mokk/config"
	"time"
)

// JsonHandler provides a Fiber Handler for rendering JSON responses
// (1) Unmarshal the incoming request body, in the event of failure respond with 400.
// (2) Find the response body for the route, including any variants.
// (3) Unmarshal the raw response body, in the event of failure respond with 500.
// (4) Respond to the client with the found response.
func JsonHandler(svr *Server, cfg config.Options, route config.Route) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body interface{}
		err := json.Unmarshal(resolveRequestBody(c.Body()), &body)

		if err != nil {
			errRes := Response{fiber.StatusBadRequest, "Error parsing request body: invalid JSON provided"}

			if cfg.TrackRequests {
				record := svr.rLog.Record(route, string(c.Body()), errRes)
				c.Append("mokk-request-id", record.Id)
			}

			return fiber.NewError(errRes.StatusCode, errRes.Response)
		}

		res := getResponse(getParamsFromCtx(c), route)

		var resBody interface{}
		err = json.Unmarshal([]byte(res.Response), &resBody)
		if err != nil {
			errRes := Response{fiber.StatusInternalServerError, fmt.Sprintf("Error rendering response: %s", err)}
			record := svr.rLog.Record(route, string(c.Body()), errRes)

			if cfg.TrackRequests {
				c.Append("mokk-request-id", record.Id)
				return fiber.NewError(errRes.StatusCode, errRes.Response)
			}
		}

		if cfg.TrackRequests {
			record := svr.rLog.Record(route, body, res)
			c.Append("mokk-request-id", record.Id)
		}

		if route.Delay > 0 {
			time.Sleep(time.Duration(route.Delay) * time.Millisecond)
		}

		if len(c.Body()) > 0 {
			if cfg.PrintRequestBody {
				str, _ := json.Marshal(body)
				fmt.Println(string(str))
			}
		}

		return c.Status(res.StatusCode).JSON(resBody)
	}
}

// getParamsFromCtx extracts all route params into a map from the Fiber context.
func getParamsFromCtx(c *fiber.Ctx) map[string]string {
	p := map[string]string{}

	for _, param := range c.Route().Params {
		p[param] = c.Params(param)
	}

	for key, value := range c.Queries() {
		p[key] = value
	}

	return p
}

// resolveRequestBody ensures a JSON body is provided for processing
func resolveRequestBody(body []byte) []byte {
	if len(body) == 0 {
		return []byte("{}")
	}

	return body
}
