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
		depositAccountUuid = uuid.New()
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
