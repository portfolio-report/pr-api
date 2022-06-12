package securities

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs/tokenize"
	"gorm.io/gorm"
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
	Tags         []string                         `json:"tags"`
}

func searchSecuritiesResponseFromDB(s db.Security) searchSecuritiesResponse {
	securityMarkets := []searchSecuritiesResponseMarket{}
	for _, mDb := range s.SecurityMarkets {
		securityMarkets = append(securityMarkets, searchSecuritiesResponseMarketFromDB(mDb))
	}

	tags := []string{}
	for _, t := range s.Tags {
		tags = append(tags, t.Name)
	}

	return searchSecuritiesResponse{
		Uuid:         strings.Replace(s.UUID.String(), "-", "", 4),
		Name:         s.Name,
		Isin:         s.Isin,
		Wkn:          s.Wkn,
		SymbolXfra:   s.SymbolXfra,
		SymbolXnas:   s.SymbolXnas,
		SymbolXnys:   s.SymbolXnys,
		SecurityType: s.SecurityType,
		Markets:      securityMarkets,
		Tags:         tags,
	}
}

type searchSecuritiesResponseMarket struct {
	MarketCode     string      `json:"marketCode"`
	Symbol         *string     `json:"symbol"`
	FirstPriceDate *model.Date `json:"firstPriceDate"`
	LastPriceDate  *model.Date `json:"lastPriceDate"`
	CurrencyCode   string      `json:"currencyCode"`
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

// SearchSecurities lists securities matching the search query
func (h *securitiesHandler) SearchSecurities(c *gin.Context) {
	searchTerm := c.Param("searchTerm")
	securityType := c.Query("securityType")

	searchTokensRaw := tokenize.SplitByWhitespace(searchTerm)
	searchTokens, searchKeyValues := tokenize.ParseKeyValue(searchTokensRaw)

	searchTerm = strings.ToUpper(strings.Join(searchTokens, " "))

	maxResults, err := strconv.Atoi(os.Getenv("SECURITIES_SEARCH_MAX_RESULTS"))
	if err != nil {
		maxResults = 10
	}

	var securities []db.Security
	query := h.DB.Preload("SecurityMarkets").Preload("Tags")

	if securityType != "" {
		query = query.Where("security_type = ?", securityType)
	}

	// Filter securities by key value pairs, if given
	for _, kv := range searchKeyValues {
		if strings.ToLower(kv[0]) == "tag" {
			var tag db.Tag
			var securityUuids []uuid.UUID

			if err := h.DB.Preload("Securities").Take(&tag, "LOWER(name) = LOWER(?)", kv[1]).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					securityUuids = []uuid.UUID{}
				} else {
					panic(err)
				}
			} else {
				securityUuids = make([]uuid.UUID, len(tag.Securities))
				for i := range tag.Securities {
					securityUuids[i] = tag.Securities[i].UUID

				}
			}

			query = query.Where("uuid IN ?", securityUuids)
			maxResults = 0
		} else {
			// ignore unknown tags
		}
	}

	if len(searchTerm) > 0 {
		query = query.Where("isin = ? OR wkn = ? OR symbol_xfra = ? OR symbol_xnas = ? OR symbol_xnys = ? OR name % ?", searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
		query = query.Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "name <-> ?", Vars: []interface{}{searchTerm}, WithoutParentheses: true},
		})
	}

	if maxResults != 0 {
		query = query.Limit(maxResults)
	}

	if err := query.Find(&securities).Error; err != nil {
		panic(err)
	}

	response := []searchSecuritiesResponse{}
	for _, s := range securities {
		response = append(response, searchSecuritiesResponseFromDB(s))
	}

	c.JSON(http.StatusOK, response)
}
