package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/richtoms/mokk/config"
	"github.com/spf13/cobra"
	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Mokk server with your config",
	Long: `
 _______         __     __
|   |   |.-----.|  |--.|  |--.
|       ||  _  ||    < |    <
|__|_|__||_____||__|__||__|__|

Start the Mokk server using the specified config file & options.
If no custom config file has been provided then Mokk will start using an example config exposing a Users API.

Example local usage:
  mokk start --config=~/code/my-app/mokk.yml --port=8080

Example docker usage:
  docker run \
    -p 8080:80 \
    --volume ~/code/my-app/mokk.yml:/app/mokk.yml \
    mokk:latest
`,
	Run: startCmdFunc,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().String("config", "./mokk.yml", "The config file to init your Mokk server")
	startCmd.Flags().Int("port", 80, "Port to host the Mokk server on")
}

// startCmdFunc is the main function for initiating a Mokk server with all the routes / config
// defined by the end user.
func startCmdFunc(cmd *cobra.Command, args []string) {
	path := cmd.Flag("config")
	cfg, err := config.LoadConfigFromFile(path.Value.String())
	if err != nil {
		panic(err)
	}

	svr := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	fmt.Println(" _______         __     __    ")
	fmt.Println("|   |   |.-----.|  |--.|  |--.")
	fmt.Println("|       ||  _  ||    < |    < ")
	fmt.Println("|__|_|__||_____||__|__||__|__|")

	tbl := table.NewWriter()
	tbl.SetStyle(table.StyleLight)
	tbl.AppendHeader(table.Row{
		"Method",
		"Path",
		"Response Code",
	})

	port := cmd.Flag("port")
	tbl.SetCaption("Mokk server hosted at: http://localhost:%s", port.Value.String())

	for _, route := range cfg.Routes {
		tbl.AppendRow(table.Row{
			route.Method,
			route.Path,
			route.StatusCode,
		})

		svr.Add(route.Method, route.Path, jsonHandler(cfg.Options, route))
	}

	fmt.Print(tbl.Render())

	fmt.Println("")
	printLog(fmt.Sprintf("[%s] Starting Mokk server...", cfg.Name))
	if err = svr.Listen(fmt.Sprintf(":%s", port.Value.String())); err != nil {
		panic(err)
	}
}

type Response struct {
	StatusCode int
	Response   string
}

// jsonHandler provides a Fiber Handler for rendering JSON responses
func jsonHandler(cfg config.Options, route config.Route) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body interface{}
		res := getResponse(c, route)

		err := json.Unmarshal([]byte(res.Response), &body)
		if err != nil {
			printLog(fmt.Sprintf("%-10.10s | %s\t %d (%s)", route.Method, route.Path, 500, utils.StatusMessage(500)))
			printLog(fmt.Sprintf("Failed to render response: %s", err))

			return fiber.ErrInternalServerError
		}

		printLog(fmt.Sprintf("%-10.10s | %s\t %d (%s)", route.Method, route.Path, res.StatusCode, utils.StatusMessage(res.StatusCode)))

		if len(c.Body()) > 0 {
			var b interface{}
			err := json.Unmarshal(c.Body(), &b)
			if err != nil {
				fmt.Println("Error parsing request body")

				return fiber.ErrBadRequest
			}

			if cfg.PrintRequestBody {
				str, _ := json.Marshal(b)
				printLog(fmt.Sprintf("%-10.10s | %s", "Body", str))
			}
		}

		return c.Status(res.StatusCode).JSON(body)
	}
}

// printLog provides a timestamped method of logging a string.
func printLog(str string) {
	fmt.Printf("%s | %s\n", time.Now().Format(time.TimeOnly), str)
}

// getResponse attempts to find the correct response based on the request
func getResponse(c *fiber.Ctx, route config.Route) Response {
	res := Response{
		StatusCode: route.StatusCode,
		Response:   route.Response,
	}

	if len(route.Variants) > 0 {
		for _, variant := range route.Variants {
			matches := make([]bool, 0)
			for key, value := range variant.Params {
				if c.Params(key) == value {
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
