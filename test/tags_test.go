package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	handlerConfig.DB.Model(&db.User{}).Where("username = 'testuser-e2e'").Update("is_admin", true)

	a := assert.New(t)

	var securityUuids [2]string

	// Create securities
	for i := 0; i < 2; i++ {
		reqBody := gin.H{
			"name": "Test name",
		}
		body, res := jsonbody[gin.H](
			api("POST", "/securities/", reqBody, &session.Token))
		a.Equal(201, res.Code)

		securityUuids[i] = body["uuid"].(string)
	}

	// Get nonexistent tag
	{
		body, res := jsonbody[gin.H](
			api("GET", "/tags/foo", nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal("foo", body["name"])
		a.Len(body["securities"], 0)
	}

	// Create tag
	{
		reqBody := gin.H{
			"securities": []gin.H{{"uuid": securityUuids[0]}},
		}
		body, res := jsonbody[gin.H](
			api("PUT", "/tags/foo", reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal("foo", body["name"])

		securities := body["securities"].([]interface{})
		a.Len(body["securities"], 1)
		security := securities[0].(map[string]interface{})
		a.Equal(securityUuids[0], security["uuid"])
	}

	// Add second security
	{
		reqBody := gin.H{
			"securities": []gin.H{{"uuid": securityUuids[0]}, {"uuid": securityUuids[1]}},
		}
		body, res := jsonbody[gin.H](
			api("PUT", "/tags/foo", reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal("foo", body["name"])

		securities := body["securities"].([]interface{})
		a.Len(body["securities"], 2)

		for i := 0; i < 2; i++ {
			security := securities[i].(map[string]interface{})
			a.Equal(securityUuids[i], security["uuid"])
		}
	}

	// Remove first security
	{
		reqBody := gin.H{
			"securities": []gin.H{{"uuid": securityUuids[1]}},
		}
		body, res := jsonbody[gin.H](
			api("PUT", "/tags/foo", reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal("foo", body["name"])

		securities := body["securities"].([]interface{})
		a.Len(body["securities"], 1)

		security := securities[0].(map[string]interface{})
		a.Equal(securityUuids[1], security["uuid"])
	}

	// Get updated tag
	{
		body, res := jsonbody[gin.H](
			api("GET", "/tags/foo", nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal("foo", body["name"])

		securities := body["securities"].([]interface{})
		a.Len(body["securities"], 1)

		security := securities[0].(map[string]interface{})
		a.Equal(securityUuids[1], security["uuid"])
	}

	// Invalid calls
	{
		res := api("PUT", "/tags/foo", nil, &session.Token)
		a.Equal(400, res.Code)

		res = api("PUT", "/tags/foo", gin.H{"securities": []gin.H{{"uuid": "11111111-1111-1111-1111-111111111111"}}}, &session.Token)
		a.Equal(400, res.Code)
	}

	// Delete tag
	{
		res := api("DELETE", "/tags/foo", nil, &session.Token)
		a.Equal(204, res.Code)
	}
}
