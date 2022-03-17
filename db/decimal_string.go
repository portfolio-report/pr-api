package db

import (
	"github.com/shopspring/decimal"
)

// DecimalString represents a decimal value stored as string
type DecimalString string

// MarshalJSON marshals DecimalString into JSON number,
// i.e. without enclosing quotes, e.g. 1.234 instead of "1.234"
func (d DecimalString) MarshalJSON() ([]byte, error) {
	return []byte(d), nil
}

// String returns underlying string
func (d *DecimalString) String() *string {
	if d == nil {
		return nil
	}
	return (*string)(d)
}

// NullDecimal converts DecimalString to NullDecimal
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
