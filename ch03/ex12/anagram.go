package anagram

import (
	"strings"
)

func isAnagram(s1, s2 string) bool {
	m1 := counterMap(s1)
	m2 := counterMap(s2)
	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}
	return true
}

func counterMap(s string) map[rune]int {
	m := make(map[rune]int)
	for _, c := range strings.ToLower(s) {
		if c != ' ' {
			m[c]++
		}
	}
	return m
}
