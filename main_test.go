package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var app http.Handler

func req(method, tarGET string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, tarGET, body)
	res := httptest.NewRecorder()
	app.ServeHTTP(res, req)
	return res
}

func TestMain(m *testing.M) {
	app = createApp()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCors(t *testing.T) {
	a := assert.New(t)

	res := req("GET", "/", nil)
	a.Equal(http.StatusOK, res.Code)
	a.Equal("*", res.Header().Get("Access-Control-Allow-Origin"))

	res = req("OPTIONS", "/", nil)
	a.Equal(http.StatusNoContent, res.Code)
	a.Equal("*", res.Header().Get("Access-Control-Allow-Origin"))
	a.Equal(0, res.Body.Len(), "Body is empty")
}

func Test404(t *testing.T) {
	a := assert.New(t)

	res := req("GET", "/does-not-exist", nil)
	a.Equal(404, res.Code, "HTTP code is 404")
	var body map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &body)
	a.Equal(404., body["statusCode"], "statusCode in JSON response is 404")
}

func TestAuth(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.url, func(t *testing.T) {
			res := req(tc.method, tc.url, nil)
			assert.Equal(t, 401, res.Code, "Forbidden without Authorization header")

			req := httptest.NewRequest(tc.method, tc.url, nil)
			req.Header.Add("Authorization", "Bearer d050be73-442e-42e2-96ab-f048527f41e2")
			res = httptest.NewRecorder()
			app.ServeHTTP(res, req)
			assert.Equal(t, 401, res.Code, "Forbidden with invalid Authorization header")
		})
	}
}
