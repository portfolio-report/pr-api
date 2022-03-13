package libs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchElementsInArrays(t *testing.T) {
	matcher := func(l int, r int) bool {
		return l == r
	}

	testCases := []struct {
		name            string
		lefts           []int
		rights          []int
		unmatchedLefts  []int
		unmatchedRights []int
		matchedRights   []int
	}{
		{
			name:  "matches numbers",
			lefts: []int{1, 2, 3, 4, 5, 6}, rights: []int{4, 5, 6, 7, 8, 9},
			unmatchedLefts: []int{1, 2, 3}, unmatchedRights: []int{7, 8, 9}, matchedRights: []int{4, 5, 6},
		},
		{
			name:  "matches only one of matching numbers on left side",
			lefts: []int{1, 2, 3, 4, 4}, rights: []int{3, 4, 5, 6},
			unmatchedLefts: []int{1, 2, 4}, unmatchedRights: []int{5, 6}, matchedRights: []int{3, 4},
		},
		{
			name:  "matches only one of matching numbers on right side",
			lefts: []int{1, 2, 3, 4}, rights: []int{3, 3, 4, 5, 6},
			unmatchedLefts: []int{1, 2}, unmatchedRights: []int{3, 5, 6}, matchedRights: []int{3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			gotUnmatchedLefts, gotUnmatchedRights, gotMatchedRights :=
				MatchElementsInArrays(tc.lefts, tc.rights, matcher)

			assertions := assert.New(t)
			assertions.ElementsMatch(gotUnmatchedLefts, tc.unmatchedLefts, "unmatchedLefts should match")
			assertions.ElementsMatch(gotUnmatchedRights, tc.unmatchedRights, "unmatchedRights should match")
			assertions.ElementsMatch(gotMatchedRights, tc.matchedRights, "matchedRights should match")
		})
	}
}
