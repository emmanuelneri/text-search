package main

import (
	"api/internal/container"
	internalHttp "api/internal/http"
	"api/internal/test"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_start(t *testing.T) {

	t.Run("should return 200 when get users and elastic return 200", func(t *testing.T) {
		r := &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			ContentLength: 2,
			Header:        http.Header(map[string][]string{"Content-Type": {"application/json"}}),
			Body:          ioutil.NopCloser(strings.NewReader(`{}`)),
		}

		port := "9999"
		startApp(t, port, r)

		client := http.Client{
			Timeout: 1 * time.Second,
		}
		res, err := client.Get(fmt.Sprintf("http://localhost:%s/users", port))
		assert.Nil(t, err)
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("should return 400 when get users and elastic return 400", func(t *testing.T) {
		r := &http.Response{
			Status:        "400 bad request",
			StatusCode:    400,
			ContentLength: 2,
			Header:        http.Header(map[string][]string{"Content-Type": {"application/json"}}),
			Body:          ioutil.NopCloser(strings.NewReader(`{}`)),
		}

		port := "8888"
		startApp(t, port, r)

		client := http.Client{
			Timeout: 1 * time.Second,
		}
		res, err := client.Get(fmt.Sprintf("http://localhost:%s/users", port))
		assert.Nil(t, err)
		assert.Equal(t, 400, res.StatusCode)

		defer res.Body.Close()
		errorResponse := internalHttp.ErrorResponse{}
		err = json.NewDecoder(res.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "400 Bad Request response error. empty response body", errorResponse.Message)
	})

	t.Run("should return 200 when scroll users and elastic return 200", func(t *testing.T) {
		r := &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			ContentLength: 2,
			Header:        http.Header(map[string][]string{"Content-Type": {"application/json"}}),
			Body:          ioutil.NopCloser(strings.NewReader(`{}`)),
		}

		port := "6666"
		startApp(t, port, r)

		client := http.Client{
			Timeout: 1 * time.Second,
		}
		res, err := client.Get(fmt.Sprintf("http://localhost:%s/users/FGluY2x1ZGVfY29udGV4dF91dWlkDXF1ZXJ5QW5kRmV0Y2gBFkJySEhCMnR5U1ZpQS16STBpRmFHMXcAAAAAAAAALxY5QW81aWx1clE4dWxnMzlFZDI2NEln/scroll", port))
		assert.Nil(t, err)
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("should return 400 when scroll users and elastic return 400", func(t *testing.T) {
		r := &http.Response{
			Status:        "400 bad request",
			StatusCode:    400,
			ContentLength: 2,
			Header:        http.Header(map[string][]string{"Content-Type": {"application/json"}}),
			Body:          ioutil.NopCloser(strings.NewReader(`{}`)),
		}

		port := "7777"
		startApp(t, port, r)

		client := http.Client{
			Timeout: 1 * time.Second,
		}
		res, err := client.Get(fmt.Sprintf("http://localhost:%s/users/FGluY2x1ZGVfY29udGV4dF91dWlkDXF1ZXJ5QW5kRmV0Y2gBFkJySEhCMnR5U1ZpQS16STBpRmFHMXcAAAAAAAAALxY5QW81aWx1clE4dWxnMzlFZDI2NEln/scroll", port))
		assert.Nil(t, err)
		assert.Equal(t, 400, res.StatusCode)

		defer res.Body.Close()
		errorResponse := internalHttp.ErrorResponse{}
		err = json.NewDecoder(res.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "400 Bad Request response error. empty response body", errorResponse.Message)
	})
}

func startApp(t *testing.T, port string, mockResponse *http.Response) {
	viper.Set("server.port", port)

	transport := test.NewMockTransport(mockResponse)
	config := elasticsearch.Config{
		Transport: transport,
	}

	es, err := elasticsearch.NewClient(config)
	assert.Nil(t, err)

	dependencyContainer := container.DependencyContainer{Es: es}

	go func() {
		start(dependencyContainer)
	}()

	ticker := time.NewTicker(10 * time.Millisecond)
	wg := test.WaitAppHealth(ticker, fmt.Sprintf("http://localhost:%s/health", port))
	wg.Wait()
	ticker.Stop()
}
