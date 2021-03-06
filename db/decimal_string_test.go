package db

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDecimalString(t *testing.T) {
	a := assert.New(t)

	d1 := DecimalString("1.23450")
	var d_nil *DecimalString = nil
	d_invalid := DecimalString("aaa")

	// MarshalJSON
	{
		b1, err := json.Marshal(d1)
		a.Nil(err)
		a.Equal("1.23450", string(b1))

		bn, err := json.Marshal(d_nil)
		a.Nil(err)
		a.Equal("null", string(bn))
	}

	// String
	{
		a.Equal("1.23450", *d1.String())
		a.Nil(d_nil.String())
	}

	// NullDecimal
	{
		nd := d1.NullDecimal()
		a.True(nd.Valid)
		a.True(decimal.NewFromFloat(1.2345).Equal(nd.Decimal))

		a.False(d_nil.NullDecimal().Valid)

		a.Panics(func() { d_invalid.NullDecimal() })
	}

}
