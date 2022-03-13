package libs

func MatchElementsInArrays[L any, R any](lefts []L, rights []R, matcher func(L, R) bool) ([]L, []R, []R) {
	unmatchedLefts := []L{}
	unmatchedRights := rights
	matchedRights := []R{}

	for _, left := range lefts {
		innerMatcher := func(right R) bool {
			return matcher(left, right)
		}

		found, posRight, elementRight := FindElementInSlice(unmatchedRights, innerMatcher)

		if found {
			unmatchedRights = remove(unmatchedRights, posRight)
			matchedRights = append(matchedRights, elementRight)
		} else {
			unmatchedLefts = append(unmatchedLefts, left)
		}
	}

	return unmatchedLefts, unmatchedRights, matchedRights
}

func FindElementInSlice[T any](arr []T, matcher func(T) bool) (found bool, pos int, element T) {
	pos = -1
	for pos, element = range arr {
		if matcher(element) {
			found = true
			return
		}
	}

	found = false
	return
}

func remove[T any](arr []T, pos int) []T {
	arr[pos] = arr[len(arr)-1]
	return arr[:len(arr)-1]
}
