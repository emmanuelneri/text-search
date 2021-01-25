package http

import (
	"api/pkg/users"
	"api/pkg/users/mocks"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	allUsers = &users.UserPaged{
		ScrollId: "123",
		Users: []users.User{
			{ID: "36e9f58f-c4e0-4c41-84a4-81ac5e608553", Name: "Tamiris Kesheh", Username: "tamiris.kesheh"},
			{ID: "e6da5140-fae5-435e-868d-9c4e4e7a465e", Name: "Elizabeth Casali", Username: "elizabethcasali"},
			{ID: "1daac1c8-6393-4649-a115-31bf7c4b8e9d", Name: "Eberson Belmiro Raimundo", Username: "ebersonbelmiroraimundo"},
			{ID: "9ec3783c-80ac-4a0c-8226-c5b15f4490d3", Name: "Dieison Sikvs", Username: "dieisonsikvs"},
			{ID: "fd383b52-8e01-4b1b-88cd-c534db56152a", Name: "Manuela Agnoletto Cortepace", Username: "manuela.agnoletto.cortepace"},
			{ID: "87edc52d-6243-463a-a47f-0d998d0970f2", Name: "Lys Dall", Username: "lys.dall"},
		},
	}

	filteredUsers = &users.UserPaged{
		ScrollId: "123",
		Users: []users.User{
			{ID: "e6da5140-fae5-435e-868d-9c4e4e7a465e", Name: "Elizabeth Casali", Username: "elizabethcasali"},
			{ID: "1daac1c8-6393-4649-a115-31bf7c4b8e9d", Name: "Eberson Belmiro Raimundo", Username: "ebersonbelmiroraimundo"},
		},
	}
)

func TestUserHandlerImpl_HandleSearch(t *testing.T) {
	path := "/users"

	allUsersBytes, err := json.Marshal(allUsers)
	assert.Nil(t, err)

	filteredUsersBytes, err := json.Marshal(filteredUsers)
	assert.Nil(t, err)

	t.Run("when search without keyword should return a list of users", func(t *testing.T) {
		mockService := new(mocks.Service)
		mockService.On("Search", mock.Anything, "").Return(allUsers, nil)
		handler := newUserHandler(mockService)

		app := fiber.New()
		app.Get(path, handler.HandleSearch())

		res, err := app.Test(httptest.NewRequest("GET", path, nil))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(allUsersBytes), string(body))

		mockService.AssertExpectations(t)
	})

	t.Run("when search with keyword should return a filtered list of users", func(t *testing.T) {
		keyword := "E"

		mockService := new(mocks.Service)
		mockService.On("Search", mock.Anything, keyword).Return(filteredUsers, nil)
		handler := newUserHandler(mockService)

		app := fiber.New()
		app.Get(path, handler.HandleSearch())

		res, err := app.Test(httptest.NewRequest("GET", path+"?search="+keyword, nil))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(filteredUsersBytes), string(body))

		mockService.AssertExpectations(t)
	})
}

func TestUserHandlerImpl_HandleScroll(t *testing.T) {
	scrollId := "123"
	path := "/users/:scrollId/scroll"

	filteredUsersBytes, err := json.Marshal(filteredUsers)
	assert.Nil(t, err)

	t.Run("when scroll without id should return not found, path not exists", func(t *testing.T) {
		mockService := new(mocks.Service)
		handler := newUserHandler(mockService)

		app := fiber.New()
		app.Get(path, handler.HandleScroll())

		req := httptest.NewRequest("GET", fmt.Sprintf("/users/%s/scroll", ""), nil)
		res, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("when scroll with id should return a list of users", func(t *testing.T) {
		mockService := new(mocks.Service)
		mockService.On("Scroll", scrollId).Return(filteredUsers, nil)
		handler := newUserHandler(mockService)

		app := fiber.New()
		app.Get(path, handler.HandleScroll())

		res, err := app.Test(httptest.NewRequest("GET", fmt.Sprintf("/users/%s/scroll", scrollId), nil))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(filteredUsersBytes), string(body))

		mockService.AssertExpectations(t)
	})
}
