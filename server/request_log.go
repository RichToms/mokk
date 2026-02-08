package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/utils/v2"
	"github.com/google/uuid"
	"github.com/richtoms/mokk/config"
	"github.com/richtoms/mokk/logging"
)

type Request struct {
	Body interface{} `json:"body"`
}

type LogEntry struct {
	Id        string       `json:"id"`
	Route     config.Route `json:"route"`
	Request   Request      `json:"request"`
	Response  Response     `json:"response"`
	Timestamp time.Time    `json:"timestamp"`
}

type RequestLog struct {
	logger  logging.Logger
	Entries map[string]LogEntry
}

func NewRequestLog(logger logging.Logger) RequestLog {
	return RequestLog{logger, map[string]LogEntry{}}
}

// Record adds a new entry to the RequestLog
func (l *RequestLog) Record(route config.Route, request interface{}, response Response) LogEntry {
	entry := LogEntry{
		Id:    uuid.NewString(),
		Route: route,
		Request: Request{
			Body: request,
		},
		Response:  response,
		Timestamp: time.Now(),
	}

	l.Entries[entry.Id] = entry

	l.printEntry(entry)

	return entry
}

// printEntry outputs to the console a timestamped line detailing the new entry.
func (l *RequestLog) printEntry(entry LogEntry) {
	l.logger.TimestampedRow(
		fmt.Sprintf("%-10.10s | %s\t %d (%s) | %s", entry.Route.Method, entry.Route.Path, entry.Response.StatusCode, utils.StatusMessage(entry.Response.StatusCode), entry.Id),
	)
}

// listEntriesHandler creates an HTTP handler for listing all entries in the RequestLog.
func listEntriesHandler(r RequestLog) fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"entries": r.Entries,
		})
	}
}

// getEntryHandler creates an HTTP handler for retrieving a single entry from the RequestLog by ID.
func getEntryHandler(r RequestLog) fiber.Handler {
	return func(c fiber.Ctx) error {
		if v, ok := r.Entries[c.Params("request")]; ok {
			return c.JSON(v)
		}

		return fiber.ErrNotFound
	}
}
