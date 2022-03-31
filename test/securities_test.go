package test

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/stretchr/testify/assert"
)

func TestSecurities(t *testing.T) {
	handlerConfig.DB.Model(&db.User{}).Where("username = 'testuser-e2e'").Update("is_admin", true)

	a := assert.New(t)

	var securityUuid string

	// Invalid create
	{
		res := api("POST", "/securities/", nil, &session.Token)
		a.Equal(400, res.Code)
	}

	// Create security
	{
		reqBody := gin.H{
			"name": "Test name",
		}
		body, res := jsonbody[gin.H](
			api("POST", "/securities/", reqBody, &session.Token))
		a.Equal(201, res.Code)
		a.Equal("Test name", body["name"])

		securityUuid = body["uuid"].(string)
	}

	// Get security (admin)
	{
		body, res := jsonbody[gin.H](
			api("GET", "/securities/"+securityUuid, nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal("Test name", body["name"])
	}

	// Get security (public)
	{
		body, res := jsonbody[gin.H](
			api("GET", "/securities/uuid/"+securityUuid, nil, nil))
		a.Equal(200, res.Code)
		a.Equal("Test name", body["name"])

		// Test UUID without dashes
		body, res = jsonbody[gin.H](
			api("GET", "/securities/uuid/"+strings.Replace(securityUuid, "-", "", 4), nil, nil))
		a.Equal(200, res.Code)
		a.Equal("Test name", body["name"])
	}

	// Search security (public)
	{
		body, res := jsonbody[[]gin.H](
			api("GET", "/securities/search/test+name", nil, nil))
		a.Equal(200, res.Code)

		found := false
		for _, s := range body {
			if s["uuid"] == strings.Replace(securityUuid, "-", "", 4) {
				found = true

				a.Equal("Test name", s["name"])
			}
		}
		a.True(found)
	}

	// Get security list (admin)
	{
		body, res := jsonbody[gin.H](
			api("GET", "/securities/?search=test+name", nil, &session.Token))
		a.Equal(200, res.Code)

		entries, ok := body["entries"].([]interface{})
		a.True(ok)
		found := false
		for _, entry := range entries {
			s, ok := entry.(map[string]interface{})
			a.True(ok)
			if s["uuid"] == securityUuid {
				found = true

				a.Equal("Test name", s["name"])
			}
		}
		a.True(found)
	}

	// Update security
	{
		reqBody := gin.H{
			"securityType": "Test type",
		}
		body, res := jsonbody[gin.H](
			api("PATCH", "/securities/"+securityUuid, reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(securityUuid, body["uuid"])
		a.Equal("Test name", body["name"])
		a.Equal("Test type", body["securityType"])

		// Invalid requests
		res = api("PATCH", "/securities/11111111-1111-1111-1111-111111111111", reqBody, &session.Token)
		a.Equal(404, res.Code)

		res = api("PATCH", "/securities/"+securityUuid, nil, &session.Token)
		a.Equal(400, res.Code)
	}

	// Delete security
	{
		body, res := jsonbody[gin.H](
			api("DELETE", "/securities/"+securityUuid, nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(securityUuid, body["uuid"])
		a.Equal("Test name", body["name"])
		a.Equal("Test type", body["securityType"])

		res = api("DELETE", "/securities/"+securityUuid, nil, &session.Token)
		a.Equal(404, res.Code)
	}
}
