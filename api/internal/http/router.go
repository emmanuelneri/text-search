package http

import (
	"api/internal/container"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// @title Users Search API
// @version 1.0
// @description This is a API to search users
// @schemes http

// @host localhost:8080
// @BasePath /
func routerStart(container container.AppContainer) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler(),
	})
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Get("/health", Live())
	app.Get("/health/ready", Ready(container.ElasticSearchClient()))

	userHandler := newUserHandler(container.UserService())
	app.Get("/users", userHandler.HandleSearch())
	app.Get("/users/:scrollId/scroll", userHandler.HandleScroll())

	return app
}
