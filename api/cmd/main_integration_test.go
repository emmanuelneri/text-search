package main

import (
	"api/internal/container"
	internalHttp "api/internal/http"
	"api/internal/test"
	"api/pkg/users"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"net/http"
	"regexp"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	ctx := context.Background()
	docker, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:      "elasticsearch:7.10.1",
			Entrypoint: nil,
			Env: map[string]string{
				"discovery.type": "single-node",
			},
			ExposedPorts: []string{"9200/tcp"},
		},
		Started: true,
	})

	assert.Nil(t, err)

	err = docker.Start(ctx)
	assert.Nil(t, err)

	defer docker.Terminate(ctx)

	p, err := docker.Ports(ctx)
	assert.Nil(t, err)
	dockerPort := p["9200/tcp"][0].HostPort

	port := "8010"
	index := "users-it"
	viper.SetDefault("server.port", port)
	viper.SetDefault("elastic.url", "http://localhost:"+dockerPort)
	viper.SetDefault("elastic.index", index)
	viper.SetDefault("elastic.debug", "true")
	viper.SetDefault("paged.size", 4)
	viper.SetDefault("paged.duration", 2*time.Minute)

	dependencyContainer := container.Setup()
	go func() {
		start(dependencyContainer)
	}()

	ticker := time.NewTicker(10 * time.Second)
	wg := test.WaitAppReady(ticker, fmt.Sprintf("http://localhost:%s/health/ready", port))
	wg.Wait()
	ticker.Stop()

	sendToElastic(t, dependencyContainer.Es, index, NewUserIndex("d41d1b73-e28a-4464-a640-abfe1d913cfd", "Nadja Elias", "nadjaelias", 90))
	sendToElastic(t, dependencyContainer.Es, index, NewUserIndex("fd3d9ebc-ed89-4a5c-ad0d-c34ecbcae2ac", "Joicy Quidicomo", "joicyquidicomo", 99))
	sendToElastic(t, dependencyContainer.Es, index, NewUserIndex("efa7c6a5-7ed7-438a-875b-1fc4184b6bf0", "JOSE Cynara Katyuce", "jose.cynara.katyuce", 50))
	sendToElastic(t, dependencyContainer.Es, index, NewUserIndex("48da4ef1-9cdf-4a18-aa93-92cdef0d6482", "Joeliton RODRIGUES", "joelitonrodrigues", 1))
	sendToElastic(t, dependencyContainer.Es, index, NewUserIndex("64b7f8c1-6264-4edf-b47a-90f8692056e0", "Emannuelly Reginaldo Fofonka", "emannuelly.reginaldo.fofonka", 0))

	t.Run("get all users without search term", func(t *testing.T) {
		usersResponse := get(t, fmt.Sprintf("http://localhost:%s/users", port))

		assert.NotEmpty(t, usersResponse.ScrollId)
		assert.Equal(t, 4, len(usersResponse.Users))

		assertUser(t, usersResponse.Users[0], "fd3d9ebc-ed89-4a5c-ad0d-c34ecbcae2ac", "Joicy Quidicomo", "joicyquidicomo")
		assertUser(t, usersResponse.Users[1], "d41d1b73-e28a-4464-a640-abfe1d913cfd", "Nadja Elias", "nadjaelias")
		assertUser(t, usersResponse.Users[2], "efa7c6a5-7ed7-438a-875b-1fc4184b6bf0", "JOSE Cynara Katyuce", "jose.cynara.katyuce")
		assertUser(t, usersResponse.Users[3], "48da4ef1-9cdf-4a18-aa93-92cdef0d6482", "Joeliton RODRIGUES", "joelitonrodrigues")
	})

	t.Run("get users with keyword jo", func(t *testing.T) {
		usersResponse := get(t, fmt.Sprintf("http://localhost:%s/users?search=jo", port))

		assert.NotEmpty(t, usersResponse.ScrollId)
		assert.Equal(t, 3, len(usersResponse.Users))

		assertUser(t, usersResponse.Users[0], "fd3d9ebc-ed89-4a5c-ad0d-c34ecbcae2ac", "Joicy Quidicomo", "joicyquidicomo")
		assertUser(t, usersResponse.Users[1], "efa7c6a5-7ed7-438a-875b-1fc4184b6bf0", "JOSE Cynara Katyuce", "jose.cynara.katyuce")
		assertUser(t, usersResponse.Users[2], "48da4ef1-9cdf-4a18-aa93-92cdef0d6482", "Joeliton RODRIGUES", "joelitonrodrigues")
	})

	t.Run("get all users without search term and scroll next page", func(t *testing.T) {
		searchResponse := get(t, fmt.Sprintf("http://localhost:%s/users", port))

		assert.NotEmpty(t, searchResponse.ScrollId)
		assert.Equal(t, 4, len(searchResponse.Users))

		scrollResponse := get(t, fmt.Sprintf("http://localhost:%s/users/%s/scroll", port, searchResponse.ScrollId))

		assert.Equal(t, 1, len(scrollResponse.Users))
		assertUser(t, scrollResponse.Users[0], "64b7f8c1-6264-4edf-b47a-90f8692056e0", "Emannuelly Reginaldo Fofonka", "emannuelly.reginaldo.fofonka")

	})

	t.Run("get scroll with invalid ID should return bad request and informe reason message", func(t *testing.T) {
		randomID, err := uuid.NewUUID()
		assert.Nil(t, err)
		errorResponse := getError(t, fmt.Sprintf("http://localhost:%s/users/%s/scroll", port, randomID.String()), http.StatusBadRequest)
		assert.Equal(t, "Cannot parse scroll id", errorResponse.Message)
	})
}
func get(t *testing.T, url string) users.UserPaged {
	httClient := http.Client{
		Timeout: 1 * time.Second,
	}

	res, err := httClient.Get(url)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	defer res.Body.Close()
	usersResponse := users.UserPaged{}
	err = json.NewDecoder(res.Body).Decode(&usersResponse)
	assert.Nil(t, err)

	return usersResponse
}

func getError(t *testing.T, url string, expectedStatus int) internalHttp.ErrorResponse {
	httClient := http.Client{
		Timeout: 1 * time.Second,
	}

	res, err := httClient.Get(url)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, res.StatusCode)

	defer res.Body.Close()
	errorResponse := internalHttp.ErrorResponse{}
	err = json.NewDecoder(res.Body).Decode(&errorResponse)
	assert.Nil(t, err)

	return errorResponse
}

func assertUser(t *testing.T, user users.User, expectedID, expectedName, expectedUsername string) {
	assert.NotNil(t, user)
	assert.Equal(t, expectedID, user.ID)
	assert.Equal(t, expectedName, user.Name)
	assert.Equal(t, expectedUsername, user.Username)
}

func sendToElastic(t *testing.T, es *elasticsearch.Client, index string, user UserIndex) {
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: user.ID,
		Body:       user.toReader(),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		assert.Nil(t, err)
		return
	}
	defer res.Body.Close()

}

type UserIndex struct {
	ID        string `json:"ID"`
	Name      string `json:"Name"`
	Username  string `json:"Username"`
	Relevance int    `json:"relevance"`
}

func NewUserIndex(ID, name, username string, relevance int) UserIndex {
	return UserIndex{
		ID:        ID,
		Name:      name,
		Username:  username,
		Relevance: relevance,
	}
}

func (u UserIndex) toReader() *bytes.Reader {
	requestByte, _ := json.Marshal(u)
	return bytes.NewReader(requestByte)
}
