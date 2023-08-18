/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
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
	Short: "Start the mokk server with your config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: startCmdFunc,
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().String("config", "./mokk.yml", "The config file to init your Mokk server")
	startCmd.Flags().Int("port", 80, "Port to host the Mokk server on")
}

func startCmdFunc(cmd *cobra.Command, args []string) {
	path := cmd.Flag("config")
	cfg, err := config.LoadConfig(path.Value.String())
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

	port := cmd.Flag("port")

	tbl := table.NewWriter()
	tbl.AppendHeader(table.Row{
		"Method",
		"Path",
		"Response Code",
	})
	tbl.SetCaption("Mokk server hosted at: http://localhost:%s", port.Value.String())

	for _, route := range cfg.Routes {
		tbl.AppendRow(table.Row{
			route.Method,
			route.Path,
			route.StatusCode,
		})
		svr.Add(route.Method, route.Path, defaultHandler(route))
	}

	tbl.SetStyle(table.StyleLight)
	fmt.Print(tbl.Render())

	fmt.Printf("\n\n%s | Starting Mokk server [%s]...\n", time.Now().Format(time.TimeOnly), cfg.Name)
	if err = svr.Listen(fmt.Sprintf(":%s", port.Value.String())); err != nil {
		panic(err)
	}
}

func defaultHandler(route config.Route) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Printf("%s | %-10.10s | %s\t", time.Now().Format(time.TimeOnly), route.Method, route.Path)
		var body map[string]interface{}
		err := json.Unmarshal([]byte(route.Response), &body)
		if err != nil {
			fmt.Printf("%d (%s)\n", 500, utils.StatusMessage(500))
			fmt.Println(fmt.Sprintf("%s | Failed to render response: %s", time.Now().Format(time.TimeOnly), err))
			return fiber.ErrInternalServerError
		}

		fmt.Printf("%d (%s)\n", route.StatusCode, utils.StatusMessage(route.StatusCode))
		return c.Status(route.StatusCode).JSON(body)
	}
}
