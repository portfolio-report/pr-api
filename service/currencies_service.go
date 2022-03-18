package service

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type currenciesService struct {
	DB *gorm.DB

	// List of all currencies incl. exchange rates
	currencies []db.Currency

	// Map to resolve currencyCode into position within currencies
	code2idx map[string]int

	// Holds conversion routes between all currencies
	nextCurrencyIndex [][]int
}

// NewCurrenciesService creates and returns currencies service
func NewCurrenciesService(db *gorm.DB) models.CurrenciesService {
	s := &currenciesService{
		DB: db,
	}

	go s.calculateCurrencyConversionRoutes()

	return s
}

// modelFromDb converts currency from database into model
func (s *currenciesService) modelFromDb(c db.Currency) *model.Currency {
	base := make([]*model.Exchangerate, len(c.ExchangeratesBase))
	for i := range c.ExchangeratesBase {
		base[i] = s.exchangerateModelFromDb(c.ExchangeratesBase[i])
	}

	quote := make([]*model.Exchangerate, len(c.ExchangeratesQuote))
	for i := range c.ExchangeratesQuote {
		quote[i] = s.exchangerateModelFromDb(c.ExchangeratesQuote[i])
	}

	return &model.Currency{
		Code:               c.Code,
		ExchangeratesBase:  base,
		ExchangeratesQuote: quote,
	}
}

// exchangerateModelFromDb converts exchange rate from database into model
func (*currenciesService) exchangerateModelFromDb(e db.Exchangerate) *model.Exchangerate {
	return &model.Exchangerate{
		BaseCurrencyCode:  e.BaseCurrencyCode,
		QuoteCurrencyCode: e.QuoteCurrencyCode,
	}
}

// GetCurrencies lists currencies with exchange rates
func (s *currenciesService) GetCurrencies() ([]*model.Currency, error) {
	var currencies []db.Currency
	if err := s.DB.
		Preload("ExchangeratesBase").Preload("ExchangeratesQuote").
		Find(&currencies).Error; err != nil {
		panic(err)
	}

	response := make([]*model.Currency, len(currencies))
	for i := range currencies {
		response[i] = s.modelFromDb(currencies[i])
	}

	return response, nil
}

// UpdateExchangeRates retrieves new prices for all exchange rates
func (s *currenciesService) UpdateExchangeRates() error {
	log.Println("Updating exchange rates...")

	var exchangeRates []db.Exchangerate
	if err := s.DB.Find(&exchangeRates).Error; err != nil {
		panic(err)
	}

	today := db.DbDate{}.FromTime(time.Now())

	for _, er := range exchangeRates {
		var latestPrice db.ExchangeratePrice

		if err := s.DB.Order("date DESC").Limit(1).Find(&latestPrice, "exchangerate_id = ?", er.ID).Error; err != nil {
			panic(err)
		}

		if latestPrice.Date.Equal(today) {
			log.Printf("Skipping exchange rate %s/%s\n", er.BaseCurrencyCode, er.QuoteCurrencyCode)
			continue
		}

		if er.BaseCurrencyCode == "EUR" {
			log.Printf("Retrieving %s/%s from ECB\n", er.BaseCurrencyCode, er.QuoteCurrencyCode)

			prices, err := s.getExchangeRatePricesEcb(er.BaseCurrencyCode, er.QuoteCurrencyCode)
			if err != nil {
				return fmt.Errorf("failed to get data from ECB for %s/%s: %w", er.BaseCurrencyCode, er.QuoteCurrencyCode, err)
			}

			newPrices := []db.ExchangeratePrice{}
			for _, p := range prices {
				if p.Time.After(latestPrice.Date.Time()) {
					newPrices = append(newPrices, db.ExchangeratePrice{
						ExchangerateID: er.ID,
						Date:           db.DbDate(p.Time),
						Value:          db.DecimalString(p.Value.String()),
					})
				}
			}

			if len(newPrices) > 0 {
				log.Printf("Adding %d new price(s)\n", len(newPrices))
				if err := s.DB.Create(&newPrices).Error; err != nil {
					panic(err)
				}
			} else {
				log.Println("No new price(s) available")
			}

		} else {
			log.Printf("No source available for %s/%s\n", er.BaseCurrencyCode, er.QuoteCurrencyCode)
		}
	}

	log.Println("Updating exchange rates finished.")
	return nil
}

type timeDecimal struct {
	Time  time.Time
	Value decimal.Decimal
}

// getExchangeRatePricesEcb downloads and parses exchange rates from ECB
func (s *currenciesService) getExchangeRatePricesEcb(baseCurrencyCode string, quoteCurrencyCode string) ([]timeDecimal, error) {
	if baseCurrencyCode != "EUR" {
		panic("unknown exchange rate")
	}

	url := "https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/" +
		strings.ToLower(quoteCurrencyCode) + ".xml"

	xml, err := xmlquery.LoadURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to load XML from %s: %w", url, err)
	}

	series := []timeDecimal{}
	for _, node := range xmlquery.Find(xml, "/CompactData/DataSet/Series/Obs") {
		timeStr := node.SelectAttr("TIME_PERIOD")
		valueStr := node.SelectAttr("OBS_VALUE")

		time, err := time.Parse("2006-01-02", timeStr)
		if err != nil {
			return nil, fmt.Errorf("while parsing %s to time.Time: %w", timeStr, err)
		}
		value, err := decimal.NewFromString(valueStr)
		if err != nil {
			return nil, fmt.Errorf("while parsing %s to decimal.Decimal: %w", valueStr, err)
		}
		series = append(series, timeDecimal{
			Time:  time,
			Value: value,
		})
	}

	return series, nil
}

// calculateCurrencyConversionRoutes calculates all conversion routes between all currencies
func (s *currenciesService) calculateCurrencyConversionRoutes() {
	log.Println("Calculating currency conversion routes")

	if err := s.DB.Preload("ExchangeratesBase").Preload("ExchangeratesQuote").Find(&s.currencies).Error; err != nil {
		panic(err)
	}

	n := len(s.currencies)

	s.code2idx = make(map[string]int, n)
	for i, c := range s.currencies {
		s.code2idx[c.Code] = i
	}

	// Get edges
	edges := []libs.Edge{}
	for _, c := range s.currencies {
		for _, er := range append(c.ExchangeratesBase, c.ExchangeratesQuote...) {
			edges = append(edges, libs.Edge{From: s.code2idx[er.BaseCurrencyCode], To: s.code2idx[er.QuoteCurrencyCode], Weight: 1})
			edges = append(edges, libs.Edge{From: s.code2idx[er.QuoteCurrencyCode], To: s.code2idx[er.BaseCurrencyCode], Weight: 1})
		}
	}

	_, s.nextCurrencyIndex = libs.FloydWarshall(n, edges)

	log.Println("Calculating currency conversion routes finished.")
}

// getConversionRoute returns list of currency codes in the conversion path between two currencies
func (s *currenciesService) getConversionRoute(
	sourceCurrencyCode string,
	targetCurrencyCode string,
) ([]string, error) {
	if _, known := s.code2idx[sourceCurrencyCode]; !known {
		return nil, fmt.Errorf("unknown currency code %s", sourceCurrencyCode)
	}
	if _, known := s.code2idx[targetCurrencyCode]; !known {
		return nil, fmt.Errorf("unknown currency code %s", targetCurrencyCode)
	}

	conversionRoute := []string{sourceCurrencyCode}

	currentCode := sourceCurrencyCode
	for currentCode != targetCurrencyCode {
		nextCurrencyIdx := s.nextCurrencyIndex[s.code2idx[currentCode]][s.code2idx[targetCurrencyCode]]
		if nextCurrencyIdx == -1 {
			return nil, fmt.Errorf("no conversion route found from currency code %s to %s", currentCode, targetCurrencyCode)
		}

		nextCurrencyCode := s.currencies[nextCurrencyIdx].Code
		conversionRoute = append(conversionRoute, nextCurrencyCode)
		currentCode = nextCurrencyCode
	}

	return conversionRoute, nil
}

// getOneExchangerateValue gets the value of the exchange rate at (or before) the given date
func (s *currenciesService) getOneExchangerateValue(baseCurrencyCode string, quoteCurrencyCode string, date time.Time) *decimal.Decimal {
	var price decimal.Decimal
	result := s.DB.Raw(`SELECT p.value `+
		`FROM exchangerates e INNER JOIN exchangerates_prices p ON e.id = p.exchangerate_id `+
		`WHERE e.base_currency_code = ? AND e.quote_currency_code = ? AND p.date <= ? `+
		`ORDER BY p.date DESC LIMIT 1`,
		baseCurrencyCode, quoteCurrencyCode, date).
		Scan(&price)
	if err := result.Error; err != nil {
		panic(err)
	}

	if result.RowsAffected == 1 {
		return &price
	}
	return nil
}

// ConvertCurrencyAmount converts amount between two currencies
func (s *currenciesService) ConvertCurrencyAmount(
	amount decimal.Decimal,
	sourceCurrencyCode string,
	targetCurrencyCode string,
	date time.Time,
) (decimal.Decimal, error) {
	route, err := s.getConversionRoute(sourceCurrencyCode, targetCurrencyCode)
	if err != nil {
		return decimal.Zero, fmt.Errorf("cannot convert amount: %w", err)
	}

	if len(route) == 1 {
		return amount, nil
	}

	for i := 0; i < len(route)-1; i++ {
		current := route[i]
		next := route[i+1]

		price := s.getOneExchangerateValue(current, next, date)
		if price != nil {
			amount = amount.Mul(*price)
		} else {
			price = s.getOneExchangerateValue(next, current, date)
			if price != nil {
				amount = amount.Div(*price)
			} else {
				return decimal.Zero, fmt.Errorf("no value of exchange rate %s/%s at %s found", current, next, date)
			}
		}
	}

	return amount, nil
}
