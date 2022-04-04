package test

import (
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPortfolios(t *testing.T) {
	a := assert.New(t)

	var portfolioId string
	var portfolioIdInt uint

	var securityUuid uuid.UUID
	var depositAccountUuid uuid.UUID

	// GET /portfolios/ -> empty
	{
		body, res := jsonbody[[]gin.H](
			api("GET", "/portfolios/", nil, &session.Token))
		a.Equal(200, res.Code)
		a.Len(body, 0)
	}

	// POST /portfolios/
	{
		reqBody := gin.H{
			"name":             "Test Portfolio",
			"note":             "Test comment",
			"baseCurrencyCode": "EUR",
		}
		body, res := jsonbody[gin.H](
			api("POST", "/portfolios/", reqBody, &session.Token))
		a.Equal(201, res.Code)
		a.Equal("Test Portfolio", body["name"])
		a.Equal("Test comment", body["note"])
		a.Equal("EUR", body["baseCurrencyCode"])

		a.IsType(1., body["id"])

		portfolioIdInt = uint(body["id"].(float64))
		portfolioId = strconv.Itoa(int(portfolioIdInt))
	}

	// GET /porfolios/
	{
		body, res := jsonbody[[]gin.H](
			api("GET", "/portfolios/", nil, &session.Token))
		a.Equal(200, res.Code)
		a.Len(body, 1)
		a.Equal("Test Portfolio", body[0]["name"])
		a.Equal("Test comment", body[0]["note"])
		a.Equal("EUR", body[0]["baseCurrencyCode"])
		a.Equal(portfolioIdInt, uint(body[0]["id"].(float64)))
	}

	// PUT /portfolios/$id
	{
		reqBody := gin.H{
			"name":             "changed name",
			"note":             "different note",
			"baseCurrencyCode": "USD",
		}
		body, res := jsonbody[gin.H](
			api("PUT", "/portfolios/"+portfolioId, reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal("changed name", body["name"])
		a.Equal("different note", body["note"])
		a.Equal("USD", body["baseCurrencyCode"])
		a.Equal(portfolioIdInt, uint(body["id"].(float64)))
	}

	// GET /portfolios/$id
	{
		body, res := jsonbody[gin.H](
			api("GET", "/portfolios/"+portfolioId, nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal("changed name", body["name"])
		a.Equal("different note", body["note"])
		a.Equal("USD", body["baseCurrencyCode"])
		a.Equal(portfolioIdInt, uint(body["id"].(float64)))
	}

	// GET /portfolios/$id/securities/ -> empty
	{
		body, res := jsonbody[[]gin.H](
			api("GET", "/portfolios/"+portfolioId+"/securities/", nil, &session.Token))
		a.Equal(200, res.Code)
		a.Len(body, 0)
	}

	// PUT /portfolios/$id/securities/$uuid -> Create
	{
		securityUuid = uuid.New()
		reqBody := gin.H{
			"name":          "Test security",
			"currencyCode":  "EUR",
			"isin":          "DE123",
			"wkn":           "123456",
			"symbol":        "S",
			"active":        true,
			"note":          "Test comment",
			"securityUuid":  nil,
			"updatedAt":     "2022-01-31T11:11:11Z",
			"calendar":      nil,
			"feed":          nil,
			"feedUrl":       nil,
			"latestFeed":    nil,
			"latestFeedUrl": nil,
			"events":        []any{},
			"properties":    nil,
		}
		body, res := jsonbody[gin.H](
			api("PUT", "/portfolios/"+portfolioId+"/securities/"+securityUuid.String(), reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(securityUuid.String(), body["uuid"])
		a.Equal("Test security", body["name"])
		a.Equal("EUR", body["currencyCode"])
		a.Equal("DE123", body["isin"])
		a.Equal("123456", body["wkn"])
		a.Equal("S", body["symbol"])
		a.Equal(true, body["active"])
		a.Equal("Test comment", body["note"])
		a.Nil(body["securityUuid"])
		a.Equal("2022-01-31T11:11:11Z", body["updatedAt"])
		a.Nil(body["calendar"])
		a.Nil(body["feed"])
		a.Nil(body["feedUrl"])
		a.Nil(body["latestFeed"])
		a.Nil(body["latestFeedUrl"])
		a.Equal([]any{}, body["events"])
		a.Equal([]any{}, body["properties"])
	}

	// GET /portfolios/$id/securities/
	{
		body, res := jsonbody[[]gin.H](
			api("GET", "/portfolios/"+portfolioId+"/securities/", nil, &session.Token))
		a.Equal(200, res.Code)
		a.Len(body, 1)
		s := body[0]
		a.Equal(securityUuid.String(), s["uuid"])
		a.Equal("Test security", s["name"])
		a.Equal("EUR", s["currencyCode"])
		a.Equal("DE123", s["isin"])
		a.Equal("123456", s["wkn"])
		a.Equal("S", s["symbol"])
		a.Equal(true, s["active"])
		a.Equal("Test comment", s["note"])
		a.Nil(s["securityUuid"])
		a.Equal("2022-01-31T11:11:11Z", s["updatedAt"])
		a.Nil(s["calendar"])
		a.Nil(s["feed"])
		a.Nil(s["feedUrl"])
		a.Nil(s["latestFeed"])
		a.Nil(s["latestFeedUrl"])
		a.Equal([]any{}, s["events"])
		a.Equal([]any{}, s["properties"])
	}

	// PUT /portfolios/$id/securities/$uuid -> Update
	{
		reqBody := gin.H{
			"name":          "changed name",
			"currencyCode":  "USD",
			"isin":          "DE456",
			"wkn":           "654321",
			"symbol":        "",
			"active":        false,
			"note":          "changed comment",
			"securityUuid":  nil,
			"updatedAt":     "2022-01-31T09:09:09Z",
			"calendar":      nil,
			"feed":          nil,
			"feedUrl":       nil,
			"latestFeed":    nil,
			"latestFeedUrl": nil,
			"events":        []any{},
			"properties":    nil,
		}
		s, res := jsonbody[gin.H](
			api("PUT", "/portfolios/"+portfolioId+"/securities/"+securityUuid.String(), reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(securityUuid.String(), s["uuid"])
		a.Equal("changed name", s["name"])
		a.Equal("USD", s["currencyCode"])
		a.Equal("DE456", s["isin"])
		a.Equal("654321", s["wkn"])
		a.Equal("", s["symbol"])
		a.Equal(false, s["active"])
		a.Equal("changed comment", s["note"])
		a.Nil(s["securityUuid"])
		a.Equal("2022-01-31T09:09:09Z", s["updatedAt"])
		a.Nil(s["calendar"])
		a.Nil(s["feed"])
		a.Nil(s["feedUrl"])
		a.Nil(s["latestFeed"])
		a.Nil(s["latestFeedUrl"])
		a.Equal([]any{}, s["events"])
		a.Equal([]any{}, s["properties"])
	}

	// GET /portfolios/$id/accounts/ -> empty
	{
		body, res := jsonbody[[]gin.H](
			api("GET", "/portfolios/"+portfolioId+"/accounts/", nil, &session.Token))
		a.Equal(200, res.Code)
		a.Len(body, 0)
	}

	// PUT /portfolios/$id/accounts/$uuid -> Create
	{
		depositAccountUuid = uuid.New()
		reqBody := gin.H{
			"type":         "deposit",
			"name":         "Test deposit",
			"note":         "Test comment",
			"currencyCode": "EUR",
			"active":       true,
			"updatedAt":    "2022-01-31T11:11:11Z",
		}
		body, res := jsonbody[gin.H](
			api("PUT", "/portfolios/"+portfolioId+"/accounts/"+depositAccountUuid.String(), reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(depositAccountUuid.String(), body["uuid"])
		a.Equal("deposit", body["type"])
		a.Equal("Test deposit", body["name"])
		a.Equal("Test comment", body["note"])
		a.Equal("EUR", body["currencyCode"])
		a.Equal(true, body["active"])
		a.Equal("2022-01-31T11:11:11Z", body["updatedAt"])
	}

	// GET /portfolios/$id/accounts/
	{
		body, res := jsonbody[[]gin.H](
			api("GET", "/portfolios/"+portfolioId+"/accounts/", nil, &session.Token))
		a.Equal(200, res.Code)
		a.Len(body, 1)
		a.Equal(depositAccountUuid.String(), body[0]["uuid"])
		a.Equal("deposit", body[0]["type"])
		a.Equal("Test deposit", body[0]["name"])
		a.Equal("Test comment", body[0]["note"])
		a.Equal("EUR", body[0]["currencyCode"])
		a.Equal(true, body[0]["active"])
		a.Equal("2022-01-31T11:11:11Z", body[0]["updatedAt"])
	}

	// PUT /portfolios/$id/accounts/$uuid -> Update
	{
		reqBody := gin.H{
			"type":         "deposit",
			"name":         "changed name",
			"note":         "changed comment",
			"currencyCode": "USD",
			"active":       false,
			"updatedAt":    "2022-01-31T09:09:09Z",
		}
		body, res := jsonbody[gin.H](
			api("PUT", "/portfolios/"+portfolioId+"/accounts/"+depositAccountUuid.String(), reqBody, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(depositAccountUuid.String(), body["uuid"])
		a.Equal("deposit", body["type"])
		a.Equal("changed name", body["name"])
		a.Equal("changed comment", body["note"])
		a.Equal("USD", body["currencyCode"])
		a.Equal(false, body["active"])
		a.Equal("2022-01-31T09:09:09Z", body["updatedAt"])
	}

	// Invalid portfolio requests
	{
		res := api("POST", "/portfolios/", nil, &session.Token)
		a.Equal(400, res.Code)

		reqBody := gin.H{
			"name":             "changed name",
			"note":             "different note",
			"baseCurrencyCode": "XXX",
		}
		res = api("POST", "/portfolios/", reqBody, &session.Token)
		a.Equal(400, res.Code)

		res = api("PUT", "/portfolios/"+portfolioId, nil, &session.Token)
		a.Equal(400, res.Code)

		res = api("PUT", "/portfolios/"+portfolioId, reqBody, &session.Token)
		a.Equal(400, res.Code)
	}

	// Invalid security requests
	{
		res := api("PUT", "/portfolios/"+portfolioId+"/securities/invalid-uuid", nil, &session.Token)
		a.Equal(404, res.Code)

		res = api("PUT", "/portfolios/"+portfolioId+"/securities/"+securityUuid.String(), nil, &session.Token)
		a.Equal(400, res.Code)

		reqBody := gin.H{
			"name":          "changed name",
			"currencyCode":  "XXX",
			"isin":          "DE456",
			"wkn":           "654321",
			"symbol":        "",
			"active":        false,
			"note":          "changed comment",
			"securityUuid":  nil,
			"updatedAt":     "2022-01-31T09:09:09Z",
			"calendar":      nil,
			"feed":          nil,
			"feedUrl":       nil,
			"latestFeed":    nil,
			"latestFeedUrl": nil,
			"events":        []any{},
			"properties":    nil,
		}
		res = api("PUT", "/portfolios/"+portfolioId+"/securities/"+securityUuid.String(), reqBody, &session.Token)
		a.Equal(400, res.Code)

		res = api("DELETE", "/portfolios/"+portfolioId+"/securities/invalid-uuid", nil, &session.Token)
		a.Equal(404, res.Code)
	}

	// Invalid account requests
	{
		res := api("PUT", "/portfolios/"+portfolioId+"/accounts/invalid-uuid", nil, &session.Token)
		a.Equal(404, res.Code)

		res = api("PUT", "/portfolios/"+portfolioId+"/accounts/"+depositAccountUuid.String(), nil, &session.Token)
		a.Equal(400, res.Code)

		reqBody := gin.H{
			"type":         "deposit",
			"name":         "name",
			"note":         "comment",
			"currencyCode": "XXX",
			"active":       true,
			"updatedAt":    "2022-01-31T11:11:11Z",
		}
		res = api("PUT", "/portfolios/"+portfolioId+"/accounts/"+depositAccountUuid.String(), reqBody, &session.Token)
		a.Equal(400, res.Code)

		res = api("DELETE", "/portfolios/"+portfolioId+"/accounts/invalid-uuid", nil, &session.Token)
		a.Equal(404, res.Code)
	}

	// DELETE /portfolios/$id/securities/$uuid
	{
		s, res := jsonbody[gin.H](
			api("DELETE", "/portfolios/"+portfolioId+"/securities/"+securityUuid.String(), nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(securityUuid.String(), s["uuid"])
		a.Equal("changed name", s["name"])
		a.Equal("USD", s["currencyCode"])
		a.Equal("DE456", s["isin"])
		a.Equal("654321", s["wkn"])
		a.Equal("", s["symbol"])
		a.Equal(false, s["active"])
		a.Equal("changed comment", s["note"])
		a.Nil(s["securityUuid"])
		a.Equal("2022-01-31T09:09:09Z", s["updatedAt"])
		a.Nil(s["calendar"])
		a.Nil(s["feed"])
		a.Nil(s["feedUrl"])
		a.Nil(s["latestFeed"])
		a.Nil(s["latestFeedUrl"])
		a.Equal([]any{}, s["events"])
		a.Equal([]any{}, s["properties"])

		res = api("DELETE", "/portfolios/"+portfolioId+"/securities/"+securityUuid.String(), nil, &session.Token)
		a.Equal(404, res.Code)
	}

	// DELETE /portfolios/$id/accounts/$uuid
	{
		body, res := jsonbody[gin.H](
			api("DELETE", "/portfolios/"+portfolioId+"/accounts/"+depositAccountUuid.String(), nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal(depositAccountUuid.String(), body["uuid"])
		a.Equal("deposit", body["type"])
		a.Equal("changed name", body["name"])
		a.Equal("changed comment", body["note"])
		a.Equal("USD", body["currencyCode"])
		a.Equal(false, body["active"])
		a.Equal("2022-01-31T09:09:09Z", body["updatedAt"])

		res = api("DELETE", "/portfolios/"+portfolioId+"/accounts/"+depositAccountUuid.String(), nil, &session.Token)
		a.Equal(404, res.Code)
	}

	// DELETE /portfolios/$id
	{
		body, res := jsonbody[gin.H](
			api("DELETE", "/portfolios/"+portfolioId, nil, &session.Token))
		a.Equal(200, res.Code)
		a.Equal("changed name", body["name"])
		a.Equal("different note", body["note"])
		a.Equal("USD", body["baseCurrencyCode"])
		a.Equal(portfolioIdInt, uint(body["id"].(float64)))
	}
}
