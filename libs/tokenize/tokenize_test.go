package tokenize

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseKeyValue(t *testing.T) {
	testCases := []struct {
		in        []string
		remaining []string
		keyValues [][2]string
	}{
		{
			in:        []string{"foo"},
			remaining: []string{"foo"},
			keyValues: [][2]string{},
		},
		{
			in:        []string{"key:value"},
			remaining: []string{},
			keyValues: [][2]string{{"key", "value"}},
		},
		{
			in:        []string{"foo", "key:value", "bar", "key:value2"},
			remaining: []string{"foo", "bar"},
			keyValues: [][2]string{{"key", "value"}, {"key", "value2"}},
		},
	}

	for _, tc := range testCases {
		t.Run("ParseKeyValue "+fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			gotRemaining, gotKV := ParseKeyValue(tc.in)
			assert.Equal(t, tc.remaining, gotRemaining)
			assert.Equal(t, tc.keyValues, gotKV)
		})
	}
}

func TestSplitByWhitespace(t *testing.T) {
	testCases := []struct {
		in  string
		out []string
	}{
		{
			in:  "foo",
			out: []string{"foo"},
		},
		{
			in:  " foo ",
			out: []string{"foo"},
		},
		{
			in:  "foo bar ",
			out: []string{"foo", "bar"},
		},
		{
			in:  "foo   bar ",
			out: []string{"foo", "bar"},
		},
		{
			in:  "Say \"Hello World\" ",
			out: []string{"Say", "Hello World"},
		},
		{
			in:  "Say \"Hello World\"!",
			out: []string{"Say", "Hello World!"},
		},
		{
			in:  "allow key:value",
			out: []string{"allow", "key:value"},
		},
		{
			in:  "allow key:\"quoted value\"",
			out: []string{"allow", "key:quoted value"},
		},
	}

	for _, tc := range testCases {
		t.Run("Split \""+tc.in+"\"", func(t *testing.T) {
			got := SplitByWhitespace(tc.in)
			require.Equal(t, tc.out, got)
		})
	}
}
