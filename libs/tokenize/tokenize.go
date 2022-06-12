package tokenize

import "strings"

func ParseKeyValue(tokens []string) ([]string, [][2]string) {
	remainingTokens := []string{}
	keyValues := [][2]string{}

	for _, t := range tokens {
		key, value, found := strings.Cut(t, ":")
		if found {
			keyValues = append(keyValues, [2]string{key, value})
		} else {
			remainingTokens = append(remainingTokens, t)
		}
	}

	return remainingTokens, keyValues
}

func SplitByWhitespace(s string) []string {
	tokens := []string{}

	sb := &strings.Builder{}
	quoted := false

	for _, r := range s {
		if r == '"' {
			quoted = !quoted
		} else if !quoted && r == ' ' {
			if sb.Len() > 0 {
				tokens = append(tokens, sb.String())
				sb.Reset()
			}
		} else {
			sb.WriteRune(r)
		}
	}

	if sb.Len() > 0 {
		tokens = append(tokens, sb.String())
	}

	return tokens
}
