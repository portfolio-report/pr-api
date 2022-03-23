package model

import (
	"fmt"
	"io"
	"strconv"
)

// PortfolioAccountType represents type of portfolio Account
type PortfolioAccountType string

const (
	PortfolioAccountTypeDeposit    PortfolioAccountType = "deposit"
	PortfolioAccountTypeSecurities PortfolioAccountType = "securities"
)

func (t PortfolioAccountType) isValid() bool {
	switch t {
	case PortfolioAccountTypeDeposit, PortfolioAccountTypeSecurities:
		return true
	}
	return false
}

// String returns underlying string
func (t PortfolioAccountType) String() string {
	return string(t)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (t *PortfolioAccountType) UnmarshalJSON(v []byte) error {
	str := string(v)
	str, err := strconv.Unquote(str)
	if err != nil {
		return fmt.Errorf("could not unquote string")
	}
	*t = PortfolioAccountType(str)
	if !t.isValid() {
		return fmt.Errorf("%s is not a valid PortfolioAccountType", str)
	}
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (t PortfolioAccountType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(t.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *PortfolioAccountType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("PortfolioAccountType must be string")
	}

	*t = PortfolioAccountType(str)
	if !t.isValid() {
		return fmt.Errorf("%s is not a valid PortfolioAccountType", str)
	}
	return nil
}
