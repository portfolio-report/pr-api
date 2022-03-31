package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/stretchr/testify/assert"
)

func TestSecurities(t *testing.T) {
	handlerConfig.DB.Model(&db.User{}).Where("username = 'testuser-e2e'").Update("is_admin", true)

	a := assert.New(t)

	var securityUuid string

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
	}

	// Delete security
	{
		body, res := jsonbody[gin.H](
			api("DELETE", "/securities/"+securityUuid, nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(securityUuid, body["uuid"])
		a.Equal("Test name", body["name"])
		a.Equal("Test type", body["securityType"])
	}
}
