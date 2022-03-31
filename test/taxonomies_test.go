package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/stretchr/testify/assert"
)

func TestTaxonomies(t *testing.T) {
	handlerConfig.DB.Model(&db.User{}).Where("username = 'testuser-e2e'").Update("is_admin", true)

	a := assert.New(t)

	var rootTaxonomyUuid string
	var secondTaxonomyUuid string

	// Create root taxonomy
	{
		reqBody := gin.H{
			"name": "Test name",
		}
		body, res := jsonbody[gin.H](
			api("POST", "/taxonomies/", reqBody, &session.Token))
		a.Equal(201, res.Code)
		a.Equal("Test name", body["name"])
		a.Nil(body["code"])

		rootTaxonomyUuid = body["uuid"].(string)
	}

	// Create second taxonomy
	{
		reqBody := gin.H{
			"name": "Second tax",
		}
		body, res := jsonbody[gin.H](
			api("POST", "/taxonomies/", reqBody, &session.Token))
		a.Equal(201, res.Code)
		a.Equal("Second tax", body["name"])
		a.Nil(body["code"])

		secondTaxonomyUuid = body["uuid"].(string)
	}

	// Make second taxonomy child of root
	{
		reqBody := gin.H{
			"name":       "Second tax",
			"parentUuid": rootTaxonomyUuid,
		}
		body, res := jsonbody[gin.H](
			api("PUT", "/taxonomies/"+secondTaxonomyUuid, reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(secondTaxonomyUuid, body["uuid"])
		a.Equal("Second tax", body["name"])
		a.Nil(body["code"])
		a.Equal(rootTaxonomyUuid, body["parentUuid"])
		a.Equal(rootTaxonomyUuid, body["rootUuid"])
	}

	// Get root taxonomy
	{
		body, res := jsonbody[gin.H](
			api("GET", "/taxonomies/"+rootTaxonomyUuid, nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(rootTaxonomyUuid, body["uuid"])
		a.Equal("Test name", body["name"])
		a.Nil(body["code"])

		a.Len(body["descendants"], 1)
	}

	// Update root taxonomy
	{
		reqBody := gin.H{
			"name": "Test name2",
			"code": "Test code",
		}
		body, res := jsonbody[gin.H](
			api("PUT", "/taxonomies/"+rootTaxonomyUuid, reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(rootTaxonomyUuid, body["uuid"])
		a.Equal("Test name2", body["name"])
		a.Equal("Test code", body["code"])
	}

	// Get all taxonomies
	{
		body, res := jsonbody[[]gin.H](
			api("GET", "/taxonomies/", nil, &session.Token))
		a.Equal(200, res.Code)

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
		reqBody := gin.H{
			"name":       "Second tax",
			"parentUuid": nil,
		}
		body, res := jsonbody[gin.H](
			api("PUT", "/taxonomies/"+secondTaxonomyUuid, reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(secondTaxonomyUuid, body["uuid"])
		a.Equal("Second tax", body["name"])
		a.Nil(body["code"])
		a.Nil(body["parentUuid"])
		a.Nil(body["rootUuid"])
	}

	// Delete root taxonomy
	{
		body, res := jsonbody[gin.H](
			api("DELETE", "/taxonomies/"+rootTaxonomyUuid, nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(rootTaxonomyUuid, body["uuid"])
		a.Equal("Test name2", body["name"])
		a.Equal("Test code", body["code"])
	}

	// Delete second taxonomy
	{
		res := api("DELETE", "/taxonomies/"+secondTaxonomyUuid, nil, &session.Token)
		a.Equal(200, res.Code)
	}
}
