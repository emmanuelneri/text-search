package conf

import (
	"api/internal/logs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

func ViperStart() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("elastic.url", "http://localhost:9200")
	viper.SetDefault("elastic.index", "users")
	viper.SetDefault("elastic.debug", "true")
	viper.SetDefault("paged.size", 15)
	viper.SetDefault("paged.duration", 2*time.Minute)

	bind("SERVER_PORT", "server.port")
	bind("ELASTIC_URL", "elastic.url")
	bind("ELASTIC_INDEX", "elastic.index")
	bind("ELASTIC_DEBUG", "elastic.debug")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))

	fields := make([]zap.Field, 0, len(viper.AllKeys()))
	for _, k := range viper.AllKeys() {
		fields = append(fields, zap.Any(k, viper.Get(k)))
	}
	logs.Logger.Info("config load", fields...)
}

func bind(envKey, keyViper string) {
	envValue := os.Getenv(envKey)
	if envValue != "" {
		viper.Set(keyViper, envValue)
	}
}
