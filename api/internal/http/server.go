package http

import (
	"api/internal/container"
	"api/internal/logs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func StartServer(container container.AppContainer) {
	app := routerStart(container)
	port := app.Listen(":" + viper.GetString("server.port"))
	if err := port; err != nil {
		logs.Logger.Panic("fail to start http server", zap.Error(err))
	}
}
