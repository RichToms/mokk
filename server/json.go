package server

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/richtoms/mokk/config"
)

// JsonHandler provides a Fiber Handler for rendering JSON responses
// (1) Unmarshal the incoming request body, in the event of failure respond with 400.
// (2) Find the response body for the route, including any variants.
// (3) Unmarshal the raw response body, in the event of failure respond with 500.
// (4) Respond to the client with the found response.
func JsonHandler(svr *Server, cfg config.Options, route config.Route) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body interface{}
		err := json.Unmarshal(c.Body(), &body)

		if err != nil {
			errRes := Response{fiber.StatusBadRequest, "Error parsing request body: invalid JSON provided"}
			record := svr.rLog.Record(route, string(c.Body()), errRes)

			c.Append("mokk-request-id", record.Id)
			return fiber.NewError(errRes.StatusCode, errRes.Response)
		}

		res := getResponse(getParamsFromCtx(c), route)

		var resBody interface{}
		err = json.Unmarshal([]byte(res.Response), &resBody)
		if err != nil {
			errRes := Response{fiber.StatusInternalServerError, fmt.Sprintf("Error rendering response: %s", err)}
			record := svr.rLog.Record(route, string(c.Body()), errRes)

			c.Append("mokk-request-id", record.Id)
			return fiber.NewError(errRes.StatusCode, errRes.Response)
		}

		record := svr.rLog.Record(route, body, res)
		c.Append("mokk-request-id", record.Id)

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

	return p
}
