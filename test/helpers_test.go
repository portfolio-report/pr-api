package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler"
	"github.com/portfolio-report/pr-api/server"
)

var app http.Handler
var handlerConfig *handler.Config
var user *model.User
var session *model.Session

func TestMain(m *testing.M) {
	godotenv.Load("../.env")

	// Set up app
	cfg, gorm := server.PrepareApp()
	handlerConfig = server.InitializeService(cfg, gorm)
	app = server.CreateApp(handlerConfig)

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

func api(method, target string, body any, token *string) *httptest.ResponseRecorder {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		bodyReader = bytes.NewReader(bodyBytes)
	} else {
		bodyReader = nil
	}

	req := httptest.NewRequest(method, target, bodyReader)
	if token != nil {
		req.Header.Add("Authorization", "Bearer "+*token)
	}
	res := httptest.NewRecorder()
	app.ServeHTTP(res, req)
	return res
}

func jsonbody[T any](res *httptest.ResponseRecorder) (T, *httptest.ResponseRecorder) {
	var body T
	err := json.Unmarshal(res.Body.Bytes(), &body)
	if err != nil {
		panic(err)
	}
	return body, res
}
