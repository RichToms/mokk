package server

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/richtoms/mokk/config"
)

// JsonHandler provides a Fiber Handler for rendering JSON responses
func JsonHandler(cfg config.Options, printFn func(str string), route config.Route) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body interface{}
		res := getResponse(getParamsFromCtx(c), route)

		err := json.Unmarshal([]byte(res.Response), &body)
		if err != nil {
			printFn(fmt.Sprintf("%-10.10s | %s\t %d (%s)", route.Method, route.Path, 500, utils.StatusMessage(500)))
			printFn(fmt.Sprintf("Failed to render response: %s", err))

			return fiber.ErrInternalServerError
		}

		printFn(fmt.Sprintf("%-10.10s | %s\t %d (%s)", route.Method, route.Path, res.StatusCode, utils.StatusMessage(res.StatusCode)))

		if len(c.Body()) > 0 {
			var b interface{}
			err := json.Unmarshal(c.Body(), &b)
			if err != nil {
				fmt.Println("Error parsing request body")

				return fiber.ErrBadRequest
			}

			if cfg.PrintRequestBody {
				str, _ := json.Marshal(b)
				printFn(fmt.Sprintf("%-10.10s | %s", "Body", str))
			}
		}

		return c.Status(res.StatusCode).JSON(body)
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
