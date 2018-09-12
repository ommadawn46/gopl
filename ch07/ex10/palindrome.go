package palindrome

import (
	"sort"
)

func IsPalindrome(s sort.Interface) bool {
	i, j := 0, s.Len()-1
	for i < j {
		if s.Less(i, j) || s.Less(j, i) {
			return false
		}
		i, j = i+1, j-1
	}
	return true
}
