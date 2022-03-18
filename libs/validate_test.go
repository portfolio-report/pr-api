package libs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDateYYYYMMDD(t *testing.T) {
	val := GetValidator()

	testCases := []struct {
		input string
		valid bool
	}{
		{"1990-02-01", true},
		{"2022-01-25", true},
		{"1990-02-0x", false},
		{"1990-2-1", false},
		{"1990-002-1", false},
		{"2000-13-01", false},
		{"2000-10-32", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			err := val.Var(tc.input, "DateYYYY-MM-DD")
			assert.Equal(t, err == nil, tc.valid)
		})
	}
}

func TestIsLaxUuid(t *testing.T) {
	val := GetValidator()

	testCases := []struct {
		input string
		valid bool
	}{
		{"da64d9cf-efee-43b3-9092-2434e9b29b17", true},
		{"da64d9cf-efee-43b3-90922434e9b29b17", false},
		{"da64d9cfefee-43b3-9092-2434e9b29b17", false},
		{"da64d9cf efee 43b3 9092 2434e9b29b17", false},
		{"da64d9cf-efee-43b3-9092-2434e9b29b17 ", false},
		{"da64d9cfefee43b390922434e9b29b17", true},
		{"za64d9cfefee43b390922434e9b29b1z", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			err := val.Var(tc.input, "LaxUuid")
			assert.Equal(t, err == nil, tc.valid)
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	val := GetValidator()

	testCases := []struct {
		input string
		valid bool
	}{
		{"short", false},
		{"longenough", true},
		{"StandarD", true},
		{"allowed-characters_.", true},
		{"w1thnumb3rs", true},
		{"contains space", false},
		{"mail@example.com", false},
		{"$pecia/ characters!", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			err := val.Var(tc.input, "ValidUsername")
			assert.Equal(t, err == nil, tc.valid)
		})
	}
}
