package test

import (
	"net/http"
	"testing"

	"github.com/portfolio-report/pr-api/db"
	"github.com/stretchr/testify/assert"
)

func TestCors(t *testing.T) {
	a := assert.New(t)

	res := api("GET", "/", nil, nil)
	a.Equal(http.StatusOK, res.Code)
	a.Equal("*", res.Header().Get("Access-Control-Allow-Origin"))

	res = api("OPTIONS", "/", nil, nil)
	a.Equal(http.StatusNoContent, res.Code)
	a.Equal("*", res.Header().Get("Access-Control-Allow-Origin"))
	a.Equal(0, res.Body.Len(), "Body is empty")
}

func Test404(t *testing.T) {
	a := assert.New(t)

	body, res := jsonbody[map[string]interface{}](
		api("GET", "/does-not-exist", nil, nil))
	a.Equal(404, res.Code, "HTTP code is 404")
	a.Equal(404., body["statusCode"], "statusCode in JSON response is 404")
}

func TestAuthRequired(t *testing.T) {
	testCases := []struct {
		method string
		url    string
	}{
		{"POST", "/auth/logout"},
		{"GET", "/auth/sessions"},
		{"GET", "/auth/users/me"},
		{"POST", "/auth/users/me/password"},
		{"DELETE", "/auth/users/me"},
		{"POST", "/portfolios/"},
		{"GET", "/portfolios/"},
		{"GET", "/portfolios/42"},
		{"GET", "/portfolios/string"},
		{"PUT", "/portfolios/42"},
		{"DELETE", "/portfolios/42"},
		{"GET", "/portfolios/42/accounts/"},
		{"PUT", "/portfolios/42/accounts/42"},
		{"DELETE", "/portfolios/42/accounts/42"},
		{"GET", "/portfolios/42/transactions/"},
		{"PUT", "/portfolios/42/transactions/42"},
		{"DELETE", "/portfolios/42/transactions/42"},
		{"GET", "/portfolios/42/securities/"},
		{"PUT", "/portfolios/42/securities/42"},
		{"DELETE", "/portfolios/42/securities/42"},
		{"GET", "/securities/"},
		{"POST", "/securities/"},
		{"GET", "/securities/42"},
		{"PATCH", "/securities/42"},
		{"DELETE", "/securities/42"},
		{"PATCH", "/securities/uuid/42/markets/42"},
		{"DELETE", "/securities/uuid/42/markets/42"},
		{"PUT", "/securities/uuid/42/taxonomies/42"},
		{"PUT", "/tags/42"},
		{"DELETE", "/tags/42"},
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.url, func(t *testing.T) {
			res := api(tc.method, tc.url, nil, nil)
			assert.Equal(t, 401, res.Code, "Forbidden without Authorization header")

			invalidToken := "d050be73-442e-42e2-96ab-f048527f41e2"
			res = api(tc.method, tc.url, nil, &invalidToken)
			assert.Equal(t, 401, res.Code, "Forbidden with invalid Authorization header")
		})
	}
}

func TestAdminAuth(t *testing.T) {
	handlerConfig.DB.Model(&db.User{}).Where("id = ?", user.ID).Update("is_admin", false)

	testCases := []struct {
		method string
		url    string
	}{
		{"GET", "/securities/"},
		{"POST", "/securities/"},
		{"GET", "/securities/42"},
		{"PATCH", "/securities/42"},
		{"DELETE", "/securities/42"},
		{"PATCH", "/securities/uuid/42/markets/42"},
		{"DELETE", "/securities/uuid/42/markets/42"},
		{"PUT", "/securities/uuid/42/taxonomies/42"},
		{"PUT", "/tags/42"},
		{"DELETE", "/tags/42"},
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.url, func(t *testing.T) {
			res := api(tc.method, tc.url, nil, &session.Token)
			assert.Equal(t, 401, res.Code, "Forbidden without admin privileges")
		})
	}
}
