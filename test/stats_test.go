package test

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/stretchr/testify/assert"
)

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
		body, res := jsonbody[[]gin.H](
			api("GET", "/stats/updates", nil, nil))
		a.Equal(http.StatusOK, res.Code)
		a.GreaterOrEqual(len(body), 1)

		found := false
		for _, b := range body {
			if b["version"] == "test-version" {
				found = true
				a.Equal(b["count"], 1.)
				_, err := time.Parse("2006-01-02T15:04:05Z07:00", b["firstUpdate"].(string))
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
		body, res := jsonbody[gin.H](
			api("GET", "/stats/updates/test-version", nil, nil))
		a.Equal(http.StatusOK, res.Code)
		a.NotNil(body["byDate"])
		a.NotNil(body["byCountry"])
	}

	var id int
	// GET all updates
	{
		body, res := jsonbody[gin.H](
			api("GET", "/stats/?version=test-version", nil, &session.Token))
		a.Equal(http.StatusOK, res.Code)
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
		res := api("DELETE", "/stats/"+strconv.Itoa(id), nil, &session.Token)
		a.Equal(http.StatusNoContent, res.Code)
		a.Equal(0, res.Body.Len())
	}
}
