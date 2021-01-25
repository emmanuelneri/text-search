package container

import (
	"api/internal/conf"
	"api/internal/logs"
	"api/pkg/users"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppContainer struct {
	userService users.Service
	es          *elasticsearch.Client
}

type DependencyContainer struct {
	Es *elasticsearch.Client
}

func Setup() DependencyContainer {
	es, err := conf.ElasticsearchClientStart()
	if err != nil {
		logs.Logger.Panic("fail elasticsearch connection", zap.Error(err))
	}

	d := DependencyContainer{}
	d.Es = es

	return d
}

func (d DependencyContainer) Start() AppContainer {
	userService := users.NewUserService(d.Es, conf.UsersIndex(),
		viper.GetInt("paged.size"), viper.GetDuration("paged.duration"))

	return AppContainer{
		userService: userService,
		es:          d.Es,
	}
}

func (c AppContainer) UserService() users.Service {
	return c.userService
}

func (c AppContainer) ElasticSearchClient() *elasticsearch.Client {
	return c.es
}
