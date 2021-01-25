package http

import (
	"api/internal/logs"
	"api/internal/search"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	InternalError []byte
)

type ErrorResponse struct {
	Message string
}

func init() {
	b, err := json.Marshal(ErrorResponse{Message: "Internal Server Error"})
	if err != nil {
		panic(err)
	}

	InternalError = b
}

func errorHandler() func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		requestID := ctx.Context().UserValue("requestid")

		if e, ok := err.(*search.InvalidQuery); ok {
			logs.Logger.Error("invalid query",
				zap.String("custom error", e.Error()),
				zap.Any("request-id", requestID),
				zap.Error(e.OriginError()))

			if e.ElasticError() != nil {
				response := convertErrResponse(e.ElasticError().Error.Reason)
				return ctx.Status(400).Send(response)
			}

			response := convertErrResponse(e.Error())
			return ctx.Status(400).Send(response)
		}

		if err != nil {
			if errors.Is(err, search.ErrSearchServiceUnavailable) {
				logs.Logger.Error("elasticsearch unavailable",
					zap.Any("request-id", requestID),
					zap.Error(err))

				response := convertErrResponse(search.ErrSearchServiceUnavailable.Error())
				return ctx.Status(400).Send(response)
			}

			logs.Logger.Error("handled error",
				zap.Any("request-id", requestID),
				zap.Error(err))

			return ctx.Status(500).Send(InternalError)
		}

		return err
	}
}

func convertErrResponse(msg string) []byte {
	b, err := json.Marshal(ErrorResponse{Message: msg})
	if err != nil {
		return InternalError
	}

	return b
}
