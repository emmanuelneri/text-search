package http

import (
	"api/internal/logs"
	"api/pkg/users"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type UserHandlerImpl struct {
	service users.Service
}

func newUserHandler(service users.Service) UserHandlerImpl {
	return UserHandlerImpl{service: service}
}

// @Tags User
// @Security Bearer
// @Summary Users search
// @Description get users
// @Accept  json
// @Produce  json
// @Param search query string false "name or username to search"
// @Success 200 {object} users.UserPaged
// @Failure 400,401,403,404 {object} ErrorResponse
// @Failure 500 {string} string
// @Router /users [get]
func (u UserHandlerImpl) HandleSearch() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		value := c.Query("search")

		result, err := u.service.Search(context.Background(), value)
		if err != nil {
			return err
		}

		err = c.JSON(result)
		if err != nil {
			return errors.Wrap(err, "fail to parse json")
		}

		logs.Logger.Info("search requested",
			zap.Any("request-id", c.Context().UserValue("requestid")),
			zap.String("request-time", time.Since(start).String()))

		return nil
	}
}

// @Tags User
// @Security Bearer
// @Summary Scroll Page users
// @Description get users scroll
// @Accept  json
// @Produce  json
// @Param scrollId path string true "Scroll Id"
// @Success 200 {object} users.UserPaged
// @Failure 400,401,403,404 {object} ErrorResponse
// @Failure 500 {string} string
// @Router /users/{scrollId}/scroll [get]
func (u UserHandlerImpl) HandleScroll() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		scrollID := c.Params("scrollId")

		result, err := u.service.Scroll(scrollID)
		if err != nil {
			return err
		}

		err = c.JSON(result)
		if err != nil {
			return errors.Wrap(err, "fail to parse json")
		}

		logs.Logger.Info("scroll requested",
			zap.Any("request-id", c.Context().UserValue("requestid")),
			zap.String("request-time", time.Since(start).String()))

		return nil
	}
}
