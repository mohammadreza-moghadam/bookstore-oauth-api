package rest

import (
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email": "email@email.com", "password", "the-password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@email.com", "the-password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid rest client response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email": "email@email.com", "password", "the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid login credentials", "status": 404, "error": "not_found"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@email.com", "the-password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {

}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email": "email@email.com", "password", "the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 1, "first_name": "mohammad.r", "last_name": "moghadam", "email": "email@email.com"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@email.com", "the-password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email": "email@email.com", "password", "the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 1, "first_name": "mohammad.r", "last_name": "moghadam", "email": "email@email.com"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@email.com", "the-password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "mohammad.r", user.FirstName)
	assert.EqualValues(t, "moghadam", user.LastName)
	assert.EqualValues(t, "email@email.com", user.Email)
}
