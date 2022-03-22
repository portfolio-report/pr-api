package model

import (
	"fmt"
	"io"
	"strconv"
)

// PortfolioTransactionType represents type of portfolio transaction
type PortfolioTransactionType string

const (
	PortfolioTransactionTypePayment            PortfolioTransactionType = "Payment"
	PortfolioTransactionTypeCurrencyTransfer   PortfolioTransactionType = "CurrencyTransfer"
	PortfolioTransactionTypeDepositInterest    PortfolioTransactionType = "DepositInterest"
	PortfolioTransactionTypeDepositFee         PortfolioTransactionType = "DepositFee"
	PortfolioTransactionTypeDepositTax         PortfolioTransactionType = "DepositTax"
	PortfolioTransactionTypeSecuritiesOrder    PortfolioTransactionType = "SecuritiesOrder"
	PortfolioTransactionTypeSecuritiesDividend PortfolioTransactionType = "SecuritiesDividend"
	PortfolioTransactionTypeSecuritiesFee      PortfolioTransactionType = "SecuritiesFee"
	PortfolioTransactionTypeSecuritiesTax      PortfolioTransactionType = "SecuritiesTax"
	PortfolioTransactionTypeSecuritiesTransfer PortfolioTransactionType = "SecuritiesTransfer"
)

func (t PortfolioTransactionType) isValid() bool {
	switch t {
	case PortfolioTransactionTypePayment,
		PortfolioTransactionTypeCurrencyTransfer,
		PortfolioTransactionTypeDepositInterest,
		PortfolioTransactionTypeDepositFee,
		PortfolioTransactionTypeDepositTax,
		PortfolioTransactionTypeSecuritiesOrder,
		PortfolioTransactionTypeSecuritiesDividend,
		PortfolioTransactionTypeSecuritiesFee,
		PortfolioTransactionTypeSecuritiesTax,
		PortfolioTransactionTypeSecuritiesTransfer:
		return true
	}
	return false
}

// String returns underlying string
func (t PortfolioTransactionType) String() string {
	return string(t)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (t *PortfolioTransactionType) UnmarshalJSON(v []byte) error {
	str := string(v)
	str, err := strconv.Unquote(str)
	if err != nil {
		return fmt.Errorf("could not unquote string")
	}
	*t = PortfolioTransactionType(str)
	if !t.isValid() {
		return fmt.Errorf("%s is not a valid PortfolioTransactionType", str)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (t PortfolioTransactionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *PortfolioTransactionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("PortfolioTransactionType must be string")
	}

	*t = PortfolioTransactionType(str)
	if !t.isValid() {
		return fmt.Errorf("%s is not a valid PortfolioTransactionType", str)
	}
	return nil
}
