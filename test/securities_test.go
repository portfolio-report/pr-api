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
	handlerConfig.DB.Create(&db.Market{Code: "TEST", Name: "Test market"})

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

	// Add market and prices
	{
		res := api("PATCH", "/securities/uuid/"+securityUuid+"/markets/TEST",
			gin.H{"currencyCode": "EUR", "symbol": "TST"}, &session.Token)
		a.Equal(200, res.Code)

		res = api("PATCH", "/securities/uuid/"+securityUuid+"/markets/TEST",
			gin.H{"prices": []gin.H{
				{"date": "2020-01-01", "close": 101.01},
				{"date": "2020-01-02", "close": 101.02},
				{"date": "2020-01-03", "close": 0},
			}}, &session.Token)
		a.Equal(200, res.Code)

		res = api("PATCH", "/securities/uuid/"+securityUuid+"/markets/TEST",
			gin.H{"prices": []gin.H{
				{"date": "2020-01-02", "close": 101.02},
				{"date": "2020-01-03", "close": 101.03},
				{"date": "2020-01-04", "close": 101.04},
				{"date": "2020-01-05", "close": 101.05},
			}}, &session.Token)
		a.Equal(200, res.Code)
	}

	// Get security (admin)
	{
		body, res := jsonbody[gin.H](
			api("GET", "/securities/"+securityUuid, nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal("Test name", body["name"])
		a.NotNil(body["markets"])
		markets := body["markets"].([]interface{})
		a.Len(markets, 1)
		market := markets[0].(map[string]interface{})
		a.Equal("TEST", market["marketCode"])
		a.Equal("EUR", market["currencyCode"])
		a.Equal("TST", market["symbol"])
		a.Equal("2020-01-01", market["firstPriceDate"])
		a.Equal("2020-01-05", market["lastPriceDate"])
		a.True(market["updatePrices"].(bool))
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
		markets := body["markets"].([]interface{})
		a.Len(markets, 1)
		market := markets[0].(map[string]interface{})
		a.Equal("TEST", market["marketCode"])
		a.Equal("EUR", market["currencyCode"])
		a.Equal("TST", market["symbol"])
		a.Equal("2020-01-01", market["firstPriceDate"])
		a.Equal("2020-01-05", market["lastPriceDate"])
	}

	// Get prices (public)
	{
		body, res := jsonbody[gin.H](
			api("GET", "/securities/uuid/"+strings.Replace(securityUuid, "-", "", 4)+"/markets/TEST?from=2020-01-01", nil, nil))
		a.Equal(200, res.Code)
		a.Equal("TEST", body["marketCode"])
		a.Equal("EUR", body["currencyCode"])
		a.Equal("TST", body["symbol"])
		a.Equal("2020-01-01", body["firstPriceDate"])
		a.Equal("2020-01-05", body["lastPriceDate"])
		prices := body["prices"].([]interface{})
		a.Len(prices, 5)
		a.Equal("2020-01-01", prices[0].(map[string]interface{})["date"])
		a.Equal(101.01, prices[0].(map[string]interface{})["close"])
		a.Equal("2020-01-02", prices[1].(map[string]interface{})["date"])
		a.Equal(101.02, prices[1].(map[string]interface{})["close"])
		a.Equal("2020-01-03", prices[2].(map[string]interface{})["date"])
		a.Equal(101.03, prices[2].(map[string]interface{})["close"])
		a.Equal("2020-01-04", prices[3].(map[string]interface{})["date"])
		a.Equal(101.04, prices[3].(map[string]interface{})["close"])
		a.Equal("2020-01-05", prices[4].(map[string]interface{})["date"])
		a.Equal(101.05, prices[4].(map[string]interface{})["close"])
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

	// Invalid calls
	{
		res := api("GET", "/securities/uuid/invalid-uuid", nil, nil)
		a.Equal(404, res.Code)

		res = api("GET", "/securities/uuid/11111111-1111-1111-1111-111111111111", nil, nil)
		a.Equal(404, res.Code)

		res = api("GET", "/securities/uuid/invalid-uuid/markets/TEST", nil, nil)
		a.Equal(404, res.Code)

		res = api("GET", "/securities/uuid/11111111-1111-1111-1111-111111111111/markets/TEST", nil, nil)
		a.Equal(404, res.Code)

		res = api("GET", "/securities/uuid/"+securityUuid+"/markets/TEST?from=not-a-date", nil, nil)
		a.Equal(400, res.Code)

		res = api("GET", "/securities/invalid-uuid", nil, &session.Token)
		a.Equal(404, res.Code)

		res = api("GET", "/securities/11111111-1111-1111-1111-111111111111", nil, &session.Token)
		a.Equal(404, res.Code)

		res = api("GET", "/securities/invalid-uuid", nil, &session.Token)
		a.Equal(404, res.Code)

		res = api("GET", "/securities/11111111-1111-1111-1111-111111111111", nil, &session.Token)
		a.Equal(404, res.Code)

		res = api("POST", "/securities/", nil, &session.Token)
		a.Equal(400, res.Code)

		res = api("PATCH", "/securities/invalid-uuid", nil, &session.Token)
		a.Equal(404, res.Code)

		res = api("PATCH", "/securities/11111111-1111-1111-1111-111111111111", gin.H{}, &session.Token)
		a.Equal(404, res.Code)

		res = api("PATCH", "/securities/"+securityUuid, nil, &session.Token)
		a.Equal(400, res.Code)

		res = api("PATCH", "/securities/invalid-uuid/markets/TEST", gin.H{}, &session.Token)
		a.Equal(404, res.Code)

		res = api("PATCH", "/securities/11111111-1111-1111-1111-111111111111/markets/TEST", gin.H{}, &session.Token)
		a.Equal(404, res.Code)
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

	handlerConfig.DB.Delete(&db.Market{Code: "TEST"})
}
