package service

import (
	"testing"
	"time"

	"github.com/portfolio-report/pr-api/db"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

type CurrenciesServiceTestSuite struct {
	suite.Suite
	db      *gorm.DB
	service *currenciesService
}

func (s *CurrenciesServiceTestSuite) SetupSuite() {
	godotenv.Load("../.env")

	var err error
	s.db, err = db.InitDb(ReadConfig().Db)
	s.Nil(err)

	service := NewCurrenciesService(s.db, false)
	var ok bool
	s.service, ok = service.(*currenciesService)
	s.True(ok)
}

func (s *CurrenciesServiceTestSuite) TearDownSuite() {
	sql, err := s.db.DB()
	s.Nil(err)
	sql.Close()
}

func TestCurrenciesService(t *testing.T) {
	suite.Run(t, new(CurrenciesServiceTestSuite))
}

func (s *CurrenciesServiceTestSuite) TestGetCurrencies() {
	currencies, err := s.service.GetCurrencies()
	s.Nil(err)
	s.Len(currencies, 35)

	foundEUR := false
	foundUSD := false
	foundAED := false
	for _, currency := range currencies {
		if currency.Code == "EUR" {
			foundEUR = true

			s.Len(currency.ExchangeratesBase, 32)
			s.Len(currency.ExchangeratesQuote, 0)
		}
		if currency.Code == "USD" {
			foundUSD = true

			s.Len(currency.ExchangeratesBase, 1)
			s.Len(currency.ExchangeratesQuote, 1)
		}
		if currency.Code == "AED" {
			foundAED = true

			s.Len(currency.ExchangeratesBase, 0)
			s.Len(currency.ExchangeratesQuote, 1)
		}
	}
	s.True(foundEUR)
	s.True(foundUSD)
	s.True(foundAED)
}

func (s *CurrenciesServiceTestSuite) TestGetExchangerate() {
	exchangerate, err := s.service.GetExchangerate("EUR", "USD")
	s.Nil(err)
	s.NotNil(exchangerate)

	_, err = s.service.GetExchangerate("USD", "EUR")
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *CurrenciesServiceTestSuite) TestGetExchangeratePrices() {
	exchangerate, err := s.service.GetExchangerate("EUR", "USD")
	s.Nil(err)
	s.NotNil(exchangerate)

	{
		prices, err := s.service.GetExchangeratePrices(exchangerate, nil)
		s.Nil(err)
		s.GreaterOrEqual(len(prices), 19)
		s.Equal("1999-01-04", prices[0].Date)
		s.Equal("1.17890000", prices[0].Value)
	}

	{
		from := "1999-01-05"
		prices, err := s.service.GetExchangeratePrices(exchangerate, &from)
		s.Nil(err)
		s.GreaterOrEqual(len(prices), 19)
		s.Equal("1999-01-05", prices[0].Date)
		s.Equal("1.17900000", prices[0].Value)
	}
}

func (s *CurrenciesServiceTestSuite) TestGetExchangerateValue() {
	{
		date, err := time.Parse("2006-01-02", "1999-01-05")
		s.Nil(err)
		value := s.service.getOneExchangerateValue("EUR", "USD", date)
		s.True(value.Equal(decimal.RequireFromString("1.179")))
	}

	{
		date, err := time.Parse("2006-01-02", "1999-01-10")
		s.Nil(err)
		value := s.service.getOneExchangerateValue("EUR", "USD", date)
		s.True(value.Equal(decimal.RequireFromString("1.1659")))
	}

	{
		date, err := time.Parse("2006-01-02", "1999-01-01")
		s.Nil(err)
		value := s.service.getOneExchangerateValue("EUR", "USD", date)
		s.Nil(value)
	}
}

func (s *CurrenciesServiceTestSuite) TestGetConversionRoute() {
	{
		route, err := s.service.getConversionRoute("EUR", "USD")
		s.Nil(err)
		s.Equal([]string{"EUR", "USD"}, route)
	}

	{
		route, err := s.service.getConversionRoute("USD", "EUR")
		s.Nil(err)
		s.Equal([]string{"USD", "EUR"}, route)
	}

	{
		route, err := s.service.getConversionRoute("EUR", "AED")
		s.Nil(err)
		s.Equal([]string{"EUR", "USD", "AED"}, route)
	}

	{
		_, err := s.service.getConversionRoute("EUR", "")
		s.Error(err)
	}

	{
		_, err := s.service.getConversionRoute("", "EUR")
		s.Error(err)
	}
}

func (s *CurrenciesServiceTestSuite) TestConvertCurrencyAmount() {
	{
		date, err := time.Parse("2006-01-02", "1999-01-05")
		s.Nil(err)
		value, err := s.service.ConvertCurrencyAmount(decimal.RequireFromString("100"), "EUR", "AED", date)
		s.Nil(err)
		s.True(value.Equal(decimal.RequireFromString("432.98775")), "amount should be 432.98775: %s", value.String())
	}

	{
		date, err := time.Parse("2006-01-02", "1999-01-10")
		s.Nil(err)
		value, err := s.service.ConvertCurrencyAmount(decimal.RequireFromString("100"), "EUR", "AED", date)
		s.Nil(err)
		s.True(value.Equal(decimal.RequireFromString("428.176775")), "amount should be 428.176775: %s", value.String())
	}

	{
		date, err := time.Parse("2006-01-02", "1999-01-01")
		s.Nil(err)
		_, err = s.service.ConvertCurrencyAmount(decimal.RequireFromString("100"), "EUR", "AED", date)
		s.Error(err)
	}

	{
		date, err := time.Parse("2006-01-02", "0001-01-01")
		s.Nil(err)
		value, err := s.service.ConvertCurrencyAmount(decimal.RequireFromString("100"), "EUR", "EUR", date)
		s.Nil(err)
		s.True(value.Equal(decimal.RequireFromString("100")), "amount should be 100: %s", value.String())
	}

	{
		date, err := time.Parse("2006-01-02", "1999-01-01")
		s.Nil(err)
		_, err = s.service.ConvertCurrencyAmount(decimal.RequireFromString("100"), "EUR", "", date)
		s.Error(err)
	}

	{
		date, err := time.Parse("2006-01-02", "1999-01-05")
		s.Nil(err)
		value, err := s.service.ConvertCurrencyAmount(decimal.RequireFromString("100"), "USD", "EUR", date)
		s.Nil(err)
		floatValue, _ := value.Float64()
		s.InDelta(84.81764207, floatValue, 0.00000001)

	}
}
