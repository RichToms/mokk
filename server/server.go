package server

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/pterm/pterm"
	"github.com/richtoms/mokk/config"
	"github.com/richtoms/mokk/logging"
)

type Server struct {
	logger  logging.Logger
	cfg     config.Config
	app     *fiber.App
	Options Options
	rLog    RequestLog
}

func NewServer(cfg config.Config, logger logging.Logger, opt Options) Server {
	app := fiber.New(fiber.Config{})

	rLog := NewRequestLog(logger)
	svr := Server{logger, cfg, app, opt, rLog}

	svr.addConfiguredRoutes()
	svr.addSystemRoutes(rLog)

	return svr
}

// addConfiguredRoutes adds any routes defined in the Mokk config to the server.
func (s *Server) addConfiguredRoutes() {
	for _, route := range s.cfg.Routes {
		s.app.Add([]string{route.Method}, route.Path, JsonHandler(s, s.cfg.Options, route))
	}
}

// addSystemRoutes adds any Mokk-specific routes to the server.
func (s *Server) addSystemRoutes(r RequestLog) {
	// RequestLog routes
	s.app.Get("_mokk/requests", listEntriesHandler(r))
	s.app.Get("_mokk/requests/:request", getEntryHandler(r))
}

// Listen is a wrapper around the Fiber App's listen where the host is derived from the server options.
func (s *Server) Listen() error {
	return s.app.Listen(s.Options.resolveHost(), fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}

// PrintConfig outputs the defined routes in table form and the final host that the server is listening on.
func (s *Server) PrintConfig() {
	tableData := pterm.TableData{
		{"Method", "Path", "Response code"},
	}

	for _, route := range s.cfg.Routes {
		pathStr := route.Path

		if len(route.Variants) > 0 {
			pathStr = fmt.Sprintf("%s (+%d variants)", pathStr, len(route.Variants))
		}

		tableData = append(tableData, []string{
			route.Method,
			pathStr,
			strconv.Itoa(route.StatusCode),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableData).Render()
	fmt.Println(fmt.Sprintf("\nMokk server listening on: http://%s:%s. Waiting for requests...", s.Options.Host, s.Options.Port))
}
