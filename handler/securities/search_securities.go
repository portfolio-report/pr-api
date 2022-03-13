package securities

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"gorm.io/gorm/clause"
)

type searchSecuritiesResponse struct {
	Uuid         string                           `json:"uuid"`
	Name         *string                          `json:"name"`
	Isin         *string                          `json:"isin"`
	Wkn          *string                          `json:"wkn"`
	SymbolXfra   *string                          `json:"symbolXfra"`
	SymbolXnas   *string                          `json:"symbolXnas"`
	SymbolXnys   *string                          `json:"symbolXnys"`
	SecurityType *string                          `json:"securityType"`
	Markets      []searchSecuritiesResponseMarket `json:"markets"`
}

func searchSecuritiesResponseFromDB(s db.Security) searchSecuritiesResponse {
	securityMarkets := []searchSecuritiesResponseMarket{}
	for _, mDb := range s.SecurityMarkets {
		securityMarkets = append(securityMarkets, searchSecuritiesResponseMarketFromDB(mDb))
	}

	return searchSecuritiesResponse{
		Uuid:         strings.Replace(s.UUID, "-", "", 4),
		Name:         s.Name,
		Isin:         s.Isin,
		Wkn:          s.Wkn,
		SymbolXfra:   s.SymbolXfra,
		SymbolXnas:   s.SymbolXnas,
		SymbolXnys:   s.SymbolXnys,
		SecurityType: s.SecurityType,
		Markets:      securityMarkets,
	}
}

type searchSecuritiesResponseMarket struct {
	MarketCode     string     `json:"marketCode"`
	Symbol         *string    `json:"symbol"`
	FirstPriceDate *db.DbDate `json:"firstPriceDate"`
	LastPriceDate  *db.DbDate `json:"lastPriceDate"`
	CurrencyCode   string     `json:"currencyCode"`
}

func searchSecuritiesResponseMarketFromDB(s db.SecurityMarket) searchSecuritiesResponseMarket {
	return searchSecuritiesResponseMarket{
		MarketCode:     s.MarketCode,
		Symbol:         s.Symbol,
		FirstPriceDate: s.FirstPriceDate,
		LastPriceDate:  s.LastPriceDate,
		CurrencyCode:   s.CurrencyCode,
	}
}

func (h *SecuritiesHandler) SearchSecurities(c *gin.Context) {
	searchTerm := strings.ToUpper(c.Param("searchTerm"))
	securityType := c.Query("securityType")

	maxResults, err := strconv.Atoi(os.Getenv("SECURITIES_SEARCH_MAX_RESULTS"))
	if err != nil {
		maxResults = 10
	}

	var securities []db.Security
	query := h.DB.Preload("SecurityMarkets")

	if securityType != "" {
		query = query.Where("security_type = ?", securityType)
	}

	query = query.Where("isin = ? OR wkn = ? OR symbol_xfra = ? OR symbol_xnas = ? OR symbol_xnys = ? OR name % ?", searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	query = query.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "name <-> ?", Vars: []interface{}{searchTerm}, WithoutParentheses: true},
	})
	if err := query.Limit(maxResults).Find(&securities).Error; err != nil {
		panic(err)
	}

	response := []searchSecuritiesResponse{}
	for _, s := range securities {
		response = append(response, searchSecuritiesResponseFromDB(s))
	}

	c.JSON(http.StatusOK, response)
}
