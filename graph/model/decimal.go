package model

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shopspring/decimal"
)

// Marshal/Unmarshal decimal.Decimal in GraphQL

// MarshalDecimalScalar writes decimal.Decimal to GraphQL, like: "1.23"
func MarshalDecimalScalar(d decimal.Decimal) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(`"` + d.String() + `"`))
	})
}

// UnmarshalDecimalScalar parses GraphQL to decimal.Decimal
func UnmarshalDecimalScalar(v interface{}) (decimal.Decimal, error) {
	switch v := v.(type) {
	case string:
		return decimal.NewFromString(v)
	case int64:
		return decimal.NewFromInt(v), nil
	case float64:
		return decimal.NewFromFloat(v), nil
	default:
		return decimal.Zero, fmt.Errorf("%T cannot be parsed to decimal.Decimal", v)
	}
}
