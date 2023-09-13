package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/richtoms/mokk/config"
	"github.com/richtoms/mokk/logging"
	"github.com/richtoms/mokk/server"
	"github.com/spf13/cobra"
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
	startCmd.Flags().String("port", config.DefaultPort, "Port to host the Mokk server on")
}

// startCmdFunc is the main function for initiating a Mokk server with all the routes / config
// defined by the end user.
func startCmdFunc(cmd *cobra.Command, args []string) {
	path := cmd.Flag("config")
	cfg, err := config.LoadConfigFromFile(path.Value.String())
	if err != nil {
		panic(err)
	}

	cfg.OverrideFromCommand(cmd)

	logger := logging.NewLogger()
	svr := createServer(cfg, logger)

	logger.PrintLogo()

	tbl := table.NewWriter()
	tbl.SetStyle(table.StyleLight)
	tbl.AppendHeader(table.Row{
		"Method",
		"Path",
		"Response Code",
	})

	tbl.SetCaption("Mokk server hosted at: http://%s:%s", svr.Options.Host, svr.Options.Port)

	for _, route := range cfg.Routes {
		tbl.AppendRow(table.Row{
			route.Method,
			route.Path,
			route.StatusCode,
		})
	}

	fmt.Print(tbl.Render())

	logger.NewLine()
	logger.TimestampedRow(fmt.Sprintf("[%s] Starting Mokk server...", cfg.Name))

	if err = svr.Listen(); err != nil {
		panic(err)
	}
}

// createServer builds a server instance with all dependencies from the main process.
func createServer(cfg config.Config, logger logging.Logger) server.Server {
	return server.NewServer(cfg, logger, server.Options{
		Port: cfg.Options.Port,
		Host: cfg.Options.Host,
	})
}
