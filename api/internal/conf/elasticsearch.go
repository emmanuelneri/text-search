package conf

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/spf13/viper"
	"os"
)

func ElasticsearchClientStart() (*elasticsearch.Client, error) {
	config := elasticsearch.Config{}
	config.Addresses = viper.GetStringSlice("elastic.url")
	if viper.GetBool("elastic.debug") {
		config.Logger = &estransport.CurlLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		}
	}

	return elasticsearch.NewClient(config)
}

func UsersIndex() string {
	return viper.GetString("elastic.index")
}
