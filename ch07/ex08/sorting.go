package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

type multiSort struct {
	tracks []*Track
	keys   []func(*Track, *Track) int
}

func (s *multiSort) AppendKey(key func(*Track, *Track) int) {
	s.keys = append(s.keys, key)
}

func (s *multiSort) Len() int {
	return len(s.tracks)
}

func (s *multiSort) Swap(i, j int) {
	s.tracks[i], s.tracks[j] = s.tracks[j], s.tracks[i]
}

func (s *multiSort) Less(i, j int) bool {
	x, y := s.tracks[i], s.tracks[j]
	for i := len(s.keys) - 1; i >= 0; i-- {
		if r := s.keys[i](x, y); r != 0 {
			return r < 0
		}
	}
	return false
}

func titleCompare(x, y *Track) int {
	if x.Title < y.Title {
		return -1
	}
	if x.Title > y.Title {
		return 1
	}
	return 0
}

func artistCompare(x, y *Track) int {
	if x.Artist < y.Artist {
		return -1
	}
	if x.Artist > y.Artist {
		return 1
	}
	return 0
}

func albumCompare(x, y *Track) int {
	if x.Album < y.Album {
		return -1
	}
	if x.Album > y.Album {
		return 1
	}
	return 0
}

func yearCompare(x, y *Track) int {
	if x.Year < y.Year {
		return -1
	}
	if x.Year > y.Year {
		return 1
	}
	return 0
}

func lengthCompare(x, y *Track) int {
	if x.Length < y.Length {
		return -1
	}
	if x.Length > y.Length {
		return 1
	}
	return 0
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush()
}

func main() {
	sortDemo()
	fmt.Println("\n")
	stableDemo()
}

func sortDemo() {
	fmt.Println("Sort")

	t := make([]*Track, len(tracks))
	copy(t, tracks)

	s := &multiSort{tracks: t}
	s.AppendKey(yearCompare)
	s.AppendKey(titleCompare)
	sort.Sort(s)

	printTracks(t)
}

func stableDemo() {
	fmt.Println("Stable")

	t := make([]*Track, len(tracks))
	copy(t, tracks)

	for _, key := range []func(*Track, *Track) int{
		yearCompare, titleCompare,
	} {
		s := &multiSort{tracks: t}
		s.AppendKey(key)
		sort.Stable(s)
	}

	printTracks(t)
}
