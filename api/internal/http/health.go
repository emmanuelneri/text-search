package http

import (
	"api/internal/logs"
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	StatusUP   = "UP"
	statusDOWN = "DOWN"
)

var (
	UP   []byte
	DOWN []byte
)

type Health struct {
	Status string `example:"UP,DOWN"`
}

func init() {
	UP = convert(StatusUP)
	DOWN = convert(statusDOWN)
}

func Live() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Send(UP)
	}
}

func Ready(es *elasticsearch.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		res, err := es.Ping()

		if err != nil {
			logs.Logger.Error("elastic ping error", zap.Error(err))
			return c.Send(DOWN)
		}

		if res.IsError() {
			defer res.Body.Close()
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(res.Body)
			if err != nil {
				logs.Logger.Error("elastic ping response not ok and fail to decode",
					zap.String("status", res.Status()),
					zap.Int("StatusCode", res.StatusCode),
					zap.Error(err))
				return c.Send(DOWN)
			}

			logs.Logger.Error("elastic ping response not ok",
				zap.String("body", buf.String()),
				zap.String("status", res.Status()),
				zap.Int("StatusCode", res.StatusCode))
			return c.Send(DOWN)
		}

		return c.Send(UP)
	}
}

func convert(status string) []byte {
	b, err := json.Marshal(Health{Status: status})
	if err != nil {
		panic(err)
	}

	return b
}
