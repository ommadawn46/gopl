package main

import (
	"fmt"
	"sort"

	"./sorting"
)

var tracks = []*sorting.Track{
	{"Go", "Delilah", "From the Roots Up", 2012, sorting.Length("3m38s")},
	{"Go", "Moby", "Moby", 1992, sorting.Length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, sorting.Length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, sorting.Length("4m24s")},
}

func main() {
	sortDemo()
	fmt.Println("\n")
	stableDemo()
}

func sortDemo() {
	fmt.Println("Sort")

	t := make([]*sorting.Track, len(tracks))
	copy(t, tracks)

	s := &sorting.MultiSort{Tracks: t}
	s.AppendKey(sorting.YearCompare)
	s.AppendKey(sorting.TitleCompare)
	sort.Sort(s)

	sorting.PrintTracks(t)
}

func stableDemo() {
	fmt.Println("Stable")

	t := make([]*sorting.Track, len(tracks))
	copy(t, tracks)

	for _, key := range []func(*sorting.Track, *sorting.Track) int{
		sorting.YearCompare, sorting.TitleCompare,
	} {
		s := &sorting.MultiSort{Tracks: t}
		s.AppendKey(key)
		sort.Stable(s)
	}

	sorting.PrintTracks(t)
}
