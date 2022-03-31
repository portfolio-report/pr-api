package test

import (
	"net/http"
	"testing"

	"github.com/portfolio-report/pr-api/db"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	handlerConfig.DB.Delete(&db.User{}, "username = 'testuser-e2e-auth'")

	a := assert.New(t)

	var token string

	// Register user
	{
		reqBody := map[string]string{"username": "testuser-e2e-auth", "password": "password"}
		body, res := jsonbody[map[string]string](
			api("POST", "/auth/register", reqBody, nil))
		a.Equal(http.StatusCreated, res.Code)
		a.Len(body["token"], 36)
		token = body["token"]
	}

	// Register existing user
	{
		reqBody := map[string]string{"username": "testuser-e2e-auth", "password": "password"}
		res := api("POST", "/auth/register", reqBody, nil)

		a.Equal(http.StatusBadRequest, res.Code)
	}

	// Get user details
	{
		res := api("GET", "/auth/users/me", nil, &token)
		a.Equal(http.StatusOK, res.Code)
	}

	// Log out user
	{
		res := api("POST", "/auth/logout", nil, &token)
		a.Equal(http.StatusNoContent, res.Code)
	}

	// Log in user
	{
		reqBody := map[string]string{"username": "testuser-e2e-auth", "password": "password"}
		body, res := jsonbody[map[string]string](api("POST", "/auth/login", reqBody, nil))
		a.Equal(http.StatusCreated, res.Code)
		a.Len(body["token"], 36)

		token = body["token"]
	}

	// Change password
	{
		reqBody := map[string]string{"oldPassword": "password", "newPassword": "better_password"}
		res := api("POST", "/auth/users/me/password", reqBody, &token)
		a.Equal(http.StatusCreated, res.Code)
	}

	// Invalid log in
	{
		reqBody := map[string]string{"username": "testuser-e2e-auth", "password": "password"}
		res := api("POST", "/auth/login", reqBody, nil)
		a.Equal(http.StatusUnauthorized, res.Code)

		reqBody = map[string]string{"username": "testuser-e2e-auth-wrong", "password": "password"}
		res = api("POST", "/auth/login", reqBody, nil)
		a.Equal(http.StatusUnauthorized, res.Code)
	}

	// List sessions
	{
		body, res := jsonbody[[]map[string]string](api("GET", "/auth/sessions", nil, &token))
		a.Equal(http.StatusOK, res.Code)

		a.Len(body, 1)
		a.Equal(token, body[0]["token"])
	}

	// Delete user
	{
		res := api("DELETE", "/auth/users/me", nil, &token)
		a.Equal(http.StatusNoContent, res.Code)
	}
}
