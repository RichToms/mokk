package server

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/richtoms/mokk/config"
)

// JsonHandler provides a Fiber Handler for rendering JSON responses
func JsonHandler(svr *Server, cfg config.Options, route config.Route) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body interface{}
		err := json.Unmarshal(c.Body(), &body)
		if err != nil {
			fmt.Println("Error parsing request body")

			return fiber.ErrBadRequest
		}

		res := getResponse(getParamsFromCtx(c), route)

		record := svr.rLog.Record(route, body, res)
		c.Append("mokk-request-id", record.Id)

		var resBody interface{}
		err = json.Unmarshal([]byte(res.Response), &resBody)
		if err != nil {
			//printFn(fmt.Sprintf("%-10.10s | %s\t %d (%s) | %s", route.Method, route.Path, 500, utils.StatusMessage(500), record.Id))
			//printFn(fmt.Sprintf("Failed to render response: %s", err))

			return fiber.ErrInternalServerError
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

	return p
}
