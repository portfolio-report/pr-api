package model

import (
	"fmt"
	"io"
	"strconv"
)

// PortfolioTransactionUnitType represents type of portfolio transaction
type PortfolioTransactionUnitType string

const (
	PortfolioTransactionUnitTypeBase PortfolioTransactionUnitType = "base"
	PortfolioTransactionUnitTypeTax  PortfolioTransactionUnitType = "tax"
	PortfolioTransactionUnitTypeFee  PortfolioTransactionUnitType = "fee"
)

func (t PortfolioTransactionUnitType) isValid() bool {
	switch t {
	case PortfolioTransactionUnitTypeBase,
		PortfolioTransactionUnitTypeTax,
		PortfolioTransactionUnitTypeFee:
		return true
	}
	return false
}

// String returns underlying string
func (t PortfolioTransactionUnitType) String() string {
	return string(t)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (t *PortfolioTransactionUnitType) UnmarshalJSON(v []byte) error {
	str := string(v)
	str, err := strconv.Unquote(str)
	if err != nil {
		return fmt.Errorf("could not unquote string")
	}
	*t = PortfolioTransactionUnitType(str)
	if !t.isValid() {
		return fmt.Errorf("%s is not a valid PortfolioTransactionUnitType", str)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (t PortfolioTransactionUnitType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *PortfolioTransactionUnitType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("PortfolioTransactionUnitType must be string")
	}

	*t = PortfolioTransactionUnitType(str)
	if !t.isValid() {
		return fmt.Errorf("%s is not a valid PortfolioTransactionUnitType", str)
	}
	return nil
}
