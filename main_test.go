package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler"
	"github.com/stretchr/testify/assert"
)

var app http.Handler
var handlerConfig *handler.Config
var user *model.User
var session *model.Session

func api(method, target string, body io.Reader, token *string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, body)
	if token != nil {
		req.Header.Add("Authorization", "Bearer "+*token)
	}
	res := httptest.NewRecorder()
	app.ServeHTTP(res, req)
	return res
}

func TestMain(m *testing.M) {
	// Set up app
	cfg, gorm := prepareApp()
	handlerConfig = initializeService(cfg, gorm)
	app = createApp(handlerConfig)

	// Prepare user
	handlerConfig.DB.Delete(&db.User{}, "username ='testuser-e2e'")
	var err error
	user, err = handlerConfig.UserService.Create("testuser-e2e")
	if err != nil {
		panic(err)
	}

	// Prepare session
	session, err = handlerConfig.SessionService.CreateSession(user, "e2e-test")
	if err != nil {
		panic(err)
	}

	// Run tests
	exitVal := m.Run()

	// Cleanup
	err = handlerConfig.UserService.Delete(user.ID)
	if err != nil {
		panic(err)
	}

	os.Exit(exitVal)
}

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

	res := api("GET", "/does-not-exist", nil, nil)
	a.Equal(404, res.Code, "HTTP code is 404")
	var body map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &body)
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
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.url, func(t *testing.T) {
			res := api(tc.method, tc.url, nil, nil)
			assert.Equal(t, 401, res.Code, "Forbidden without Authorization header")

			req := httptest.NewRequest(tc.method, tc.url, nil)
			req.Header.Add("Authorization", "Bearer d050be73-442e-42e2-96ab-f048527f41e2")
			res = httptest.NewRecorder()
			app.ServeHTTP(res, req)
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
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.url, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			req.Header.Add("Authorization", "Bearer "+session.Token)
			res := httptest.NewRecorder()
			app.ServeHTTP(res, req)

			assert.Equal(t, 401, res.Code, "Forbidden without admin privileges")
		})
	}
}

func TestSecurities(t *testing.T) {
	handlerConfig.DB.Model(&db.User{}).Where("username = 'testuser-e2e'").Update("is_admin", true)

	a := assert.New(t)

	var securityUuid string

	// Create security
	{
		reqBody, err := json.Marshal(gin.H{
			"name": "Test name",
		})
		a.Nil(err)
		req := httptest.NewRequest("POST", "/securities/", bytes.NewReader(reqBody))
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(201, res.Code)

		var body gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Equal("Test name", body["name"])

		securityUuid = body["uuid"].(string)
	}

	// Update security
	{
		reqBody, err := json.Marshal(gin.H{
			"securityType": "Test type",
		})
		a.Nil(err)
		req := httptest.NewRequest("PATCH", "/securities/"+securityUuid, bytes.NewReader(reqBody))
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(200, res.Code)

		var body gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Equal(securityUuid, body["uuid"])
		a.Equal("Test name", body["name"])
		a.Equal("Test type", body["securityType"])
	}

	// Delete security
	{
		req := httptest.NewRequest("DELETE", "/securities/"+securityUuid, nil)
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(200, res.Code)

		var body gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Equal(securityUuid, body["uuid"])
		a.Equal("Test name", body["name"])
		a.Equal("Test type", body["securityType"])
	}
}

func TestTaxonomies(t *testing.T) {
	handlerConfig.DB.Model(&db.User{}).Where("username = 'testuser-e2e'").Update("is_admin", true)

	a := assert.New(t)

	var rootTaxonomyUuid string
	var secondTaxonomyUuid string

	// Create root taxonomy
	{
		reqBody, err := json.Marshal(gin.H{
			"name": "Test name",
		})
		a.Nil(err)
		req := httptest.NewRequest("POST", "/taxonomies/", bytes.NewReader(reqBody))
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(201, res.Code)

		var body gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Equal("Test name", body["name"])
		a.Nil(body["code"])

		rootTaxonomyUuid = body["uuid"].(string)
	}

	// Create second taxonomy
	{
		reqBody, err := json.Marshal(gin.H{
			"name": "Second tax",
		})
		a.Nil(err)
		req := httptest.NewRequest("POST", "/taxonomies/", bytes.NewReader(reqBody))
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(201, res.Code)

		var body gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Equal("Second tax", body["name"])
		a.Nil(body["code"])

		secondTaxonomyUuid = body["uuid"].(string)
	}

	// Make second taxonomy child of root
	{
		reqBody, err := json.Marshal(gin.H{
			"name":       "Second tax",
			"parentUuid": rootTaxonomyUuid,
		})
		a.Nil(err)
		req := httptest.NewRequest("PUT", "/taxonomies/"+secondTaxonomyUuid, bytes.NewReader(reqBody))
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(200, res.Code)

		var body gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Equal(secondTaxonomyUuid, body["uuid"])
		a.Equal("Second tax", body["name"])
		a.Nil(body["code"])
		a.Equal(rootTaxonomyUuid, body["parentUuid"])
		a.Equal(rootTaxonomyUuid, body["rootUuid"])
	}

	// Get root taxonomy
	{
		req := httptest.NewRequest("GET", "/taxonomies/"+rootTaxonomyUuid, nil)
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(200, res.Code)

		var body gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Equal(rootTaxonomyUuid, body["uuid"])
		a.Equal("Test name", body["name"])
		a.Nil(body["code"])

		a.Len(body["descendants"], 1)
	}

	// Update root taxonomy
	{
		reqBody, err := json.Marshal(gin.H{
			"name": "Test name2",
			"code": "Test code",
		})
		a.Nil(err)
		req := httptest.NewRequest("PUT", "/taxonomies/"+rootTaxonomyUuid, bytes.NewReader(reqBody))
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(200, res.Code)

		var body gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Equal(rootTaxonomyUuid, body["uuid"])
		a.Equal("Test name2", body["name"])
		a.Equal("Test code", body["code"])
	}

	// Get all taxonomies
	{
		req := httptest.NewRequest("GET", "/taxonomies/", nil)
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(200, res.Code)

		var body []gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		rootFound := false
		secondFound := false
		for _, tax := range body {
			switch tax["uuid"] {
			case rootTaxonomyUuid:
				rootFound = true
			case secondTaxonomyUuid:
				secondFound = true
			}
		}
		a.True(rootFound)
		a.True(secondFound)
	}

	// Move second taxonomy out of root
	{
		reqBody, err := json.Marshal(gin.H{
			"name":       "Second tax",
			"parentUuid": nil,
		})
		a.Nil(err)
		req := httptest.NewRequest("PUT", "/taxonomies/"+secondTaxonomyUuid, bytes.NewReader(reqBody))
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(200, res.Code)

		var body gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Equal(secondTaxonomyUuid, body["uuid"])
		a.Equal("Second tax", body["name"])
		a.Nil(body["code"])
		a.Nil(body["parentUuid"])
		a.Nil(body["rootUuid"])
	}

	// Delete root taxonomy
	{
		req := httptest.NewRequest("DELETE", "/taxonomies/"+rootTaxonomyUuid, nil)
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(200, res.Code)

		var body gin.H
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Equal(rootTaxonomyUuid, body["uuid"])
		a.Equal("Test name2", body["name"])
		a.Equal("Test code", body["code"])
	}

	// Delete second taxonomy
	{
		req := httptest.NewRequest("DELETE", "/taxonomies/"+secondTaxonomyUuid, nil)
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(200, res.Code)
	}
}

func TestStats(t *testing.T) {
	handlerConfig.DB.Model(&db.User{}).Where("username = 'testuser-e2e'").Update("is_admin", true)
	handlerConfig.DB.Delete(&db.Clientupdate{}, "version = 'test-version'")

	a := assert.New(t)

	// HEAD request
	{
		res := api("HEAD", "/stats/update/name.abuchen.portfolio/test-version", nil, nil)
		a.Equal(http.StatusOK, res.Code)
		a.Equal(res.Body.Len(), 0)
	}

	// GET /updates
	{
		res := api("GET", "/stats/updates", nil, nil)
		a.Equal(http.StatusOK, res.Code)
		var body []gin.H
		err := json.Unmarshal(res.Body.Bytes(), &body)
		a.Nil(err)
		a.GreaterOrEqual(len(body), 1)

		found := false
		for _, b := range body {
			if b["version"] == "test-version" {
				found = true
				a.Equal(b["count"], 1.)
				_, err = time.Parse("2006-01-02T15:04:05Z07:00", b["firstUpdate"].(string))
				a.Nil(err)
				timestamp, err := time.Parse("2006-01-02T15:04:05Z07:00", b["lastUpdate"].(string))
				a.Nil(err)
				age := time.Now().Sub(timestamp)
				a.Positive(age)
				a.Less(age, 1*time.Second)
			}
		}
		a.True(found)
	}

	// GET /update/test-version
	{
		res := api("GET", "/stats/updates/test-version", nil, nil)
		a.Equal(http.StatusOK, res.Code)
		var body gin.H
		err := json.Unmarshal(res.Body.Bytes(), &body)
		a.Nil(err)
		a.NotNil(body["byDate"])
		a.NotNil(body["byCountry"])
	}

	var id int
	// GET all updates
	{
		req := httptest.NewRequest("GET", "/stats/?version=test-version", nil)
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(http.StatusOK, res.Code)
		var body gin.H
		err := json.Unmarshal(res.Body.Bytes(), &body)
		a.Nil(err)
		a.NotNil(body["entries"])
		a.NotNil(body["params"])

		entries, ok := body["entries"].([]interface{})
		a.True(ok)
		a.Len(entries, 1)

		entry, ok := entries[0].(map[string]interface{})
		a.True(ok)
		a.Equal("test-version", entry["version"])
		a.True(entry["country"] == "-" || entry["country"] == "")
		a.Equal("", entry["useragent"])
		id = int(entry["id"].(float64))

		timestamp, err := time.Parse("2006-01-02T15:04:05Z07:00", entry["timestamp"].(string))
		a.Nil(err)
		age := time.Now().Sub(timestamp)
		a.Positive(age)
		a.Less(age, 1*time.Second)
	}

	// DELETE update
	{
		req := httptest.NewRequest("DELETE", "/stats/"+strconv.Itoa(id), nil)
		req.Header.Add("Authorization", "Bearer "+session.Token)
		res := httptest.NewRecorder()
		app.ServeHTTP(res, req)

		a.Equal(http.StatusNoContent, res.Code)
		a.Equal(0, res.Body.Len())
	}
}

func TestAuth(t *testing.T) {
	handlerConfig.DB.Delete(&db.User{}, "username = 'testuser-e2e-auth'")

	a := assert.New(t)

	var token string

	// Register user
	{
		reqBody, err := json.Marshal(map[string]string{"username": "testuser-e2e-auth", "password": "password"})
		a.Nil(err)
		res := api("POST", "/auth/register", bytes.NewReader(reqBody), nil)

		a.Equal(http.StatusCreated, res.Code)
		var body map[string]string
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Len(body["token"], 36)
		token = body["token"]
	}

	// Register existing user
	{
		reqBody, err := json.Marshal(map[string]string{"username": "testuser-e2e-auth", "password": "password"})
		a.Nil(err)
		res := api("POST", "/auth/register", bytes.NewReader(reqBody), nil)

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
		reqBody, err := json.Marshal(map[string]string{"username": "testuser-e2e-auth", "password": "password"})
		a.Nil(err)
		res := api("POST", "/auth/login", bytes.NewReader(reqBody), nil)
		a.Equal(http.StatusCreated, res.Code)

		var body map[string]string
		json.Unmarshal(res.Body.Bytes(), &body)
		a.Len(body["token"], 36)

		token = body["token"]
	}

	// Change password
	{
		reqBody, err := json.Marshal(map[string]string{"oldPassword": "password", "newPassword": "better_password"})
		a.Nil(err)
		res := api("POST", "/auth/users/me/password", bytes.NewReader(reqBody), &token)
		a.Equal(http.StatusCreated, res.Code)
	}

	// Invalid log in
	{
		reqBody, err := json.Marshal(map[string]string{"username": "testuser-e2e-auth", "password": "password"})
		a.Nil(err)
		res := api("POST", "/auth/login", bytes.NewReader(reqBody), nil)
		a.Equal(http.StatusUnauthorized, res.Code)

		reqBody, err = json.Marshal(map[string]string{"username": "testuser-e2e-auth-wrong", "password": "password"})
		a.Nil(err)
		res = api("POST", "/auth/login", bytes.NewReader(reqBody), nil)
		a.Equal(http.StatusUnauthorized, res.Code)
	}

	// List sessions
	{
		res := api("GET", "/auth/sessions", nil, &token)
		a.Equal(http.StatusOK, res.Code)

		var body []map[string]string
		json.Unmarshal(res.Body.Bytes(), &body)

		a.Len(body, 1)
		a.Equal(token, body[0]["token"])
	}

	// Delete user
	{
		res := api("DELETE", "/auth/users/me", nil, &token)
		a.Equal(http.StatusNoContent, res.Code)
	}

}
