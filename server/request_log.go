package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	"github.com/richtoms/mokk/config"
	"github.com/richtoms/mokk/logging"
	"time"
)

type Request struct {
	Body interface{}
}

type LogEntry struct {
	Id        string
	Route     config.Route
	Request   Request
	Response  Response
	Timestamp time.Time
}

type RequestLog struct {
	logger  logging.Logger
	Entries map[string]LogEntry
}

func NewRequestLog(logger logging.Logger) RequestLog {
	return RequestLog{logger, map[string]LogEntry{}}
}

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

func (l *RequestLog) printEntry(entry LogEntry) {
	l.logger.TimestampedRow(
		fmt.Sprintf("%-10.10s | %s\t %d (%s) | %s", entry.Route.Method, entry.Route.Path, entry.Response.StatusCode, utils.StatusMessage(entry.Response.StatusCode), entry.Id),
	)
}

func listEntriesHandler(r RequestLog) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"entries": r.Entries,
		})
	}
}

func getEntryHandler(r RequestLog) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if v, ok := r.Entries[c.Params("request")]; ok {
			return c.JSON(v)
		}

		return fiber.ErrNotFound
	}
}
