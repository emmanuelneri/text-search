package main

import (
	"api/internal/conf"
	"api/internal/container"
	"api/internal/http"
	"api/internal/logs"
)

func main() {
	logs.Logger.Info("search API started")
	conf.ViperStart()

	dependencyContainer := container.Setup()
	start(dependencyContainer)
}

func start(dependencyContainer container.DependencyContainer) {
	c := dependencyContainer.Start()
	http.StartServer(c)
}
