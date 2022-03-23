package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

// Marshal/Unmarshal uuid.UUID in GraphQL

// MarshalUUIDScalar writes uuid.UUID to GraphQL
func MarshalUUIDScalar(u uuid.UUID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprint(w, strconv.Quote(u.String()))
	})
}

// UnmarshalUUIDScalar parses GraphQL to uuid.UUID
func UnmarshalUUIDScalar(v interface{}) (uuid.UUID, error) {
	s, ok := v.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("%T cannot be parsed to uuid.UUID", v)
	}
	return uuid.Parse(s)
}
