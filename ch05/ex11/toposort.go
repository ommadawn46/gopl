package main

import (
	"fmt"
	"log"
	"sort"
)

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organiczation",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organiczation"},
	"programming languages": {"data structures", "computer organiczation"},
}

var prereqs_loop = map[string][]string{
	"algorithms":     {"data structures"},
	"calculus":       {"linear algebra"},
	"linear algebra": {"calculus"}, // loop
	"compilers": {
		"data structures",
		"formal languages",
		"computer organiczation",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organiczation"},
	"programming languages": {"data structures", "computer organiczation"},
}

func main() {
	fmt.Println("[+] prereqs")
	if order, err := topoSort(prereqs); err != nil {
		log.Fatal(err)
	} else {
		for i, course := range order {
			fmt.Printf("%d:\t%s\n", i+1, course)
		}
	}
	fmt.Println("\n[+] prereqs_loop")
	if order, err := topoSort(prereqs_loop); err != nil {
		log.Fatal(err)
	} else {
		for i, course := range order {
			fmt.Printf("%d:\t%s\n", i+1, course)
		}
	}
}

func topoSort(m map[string][]string) ([]string, error) {
	var order []string
	seen := make(map[string]bool)
	var visitAll func([]string, []string) error

	contains := func(s []string, t string) bool {
		for _, a := range s {
			if a == t {
				return true
			}
		}
		return false
	}

	visitAll = func(items []string, parents []string) error {
		for _, item := range items {
			if contains(parents, item) {
				return fmt.Errorf("loop detected")
			}
			if !seen[item] {
				seen[item] = true
				if err := visitAll(m[item], append(parents, item)); err != nil {
					return err
				}
				order = append(order, item)
			}
		}
		return nil
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	err := visitAll(keys, nil)
	return order, err
}
