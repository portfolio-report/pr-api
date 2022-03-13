package db

import (
	"github.com/shopspring/decimal"
)

type DecimalString string

// DecimalString is marshalled to JSON without enclosing quotes
// e.g. 1.234 instead of "1.234"
func (d DecimalString) MarshalJSON() ([]byte, error) {
	return []byte(d), nil
}

func (d *DecimalString) String() *string {
	if d == nil {
		return nil
	}
	return (*string)(d)
}

func (d *DecimalString) NullDecimal() decimal.NullDecimal {
	if d == nil {
		return decimal.NullDecimal{Valid: false}
	}

	dec, err := decimal.NewFromString((string)(*d))
	if err != nil {
		panic(err)
	}
	return decimal.NewNullDecimal(dec)
}

func (d *DecimalString) Decimal() *decimal.Decimal {
	if d == nil {
		return nil
	}

	ret, err := decimal.NewFromString((string)(*d))
	if err != nil {
		panic(err)
	}
	return &ret
}
