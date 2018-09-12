package sorting

import (
	"time"
)

func Length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

type MultiSort struct {
	Tracks []*Track
	keys   []func(*Track, *Track) int
}

func (s *MultiSort) AppendKey(key func(*Track, *Track) int) {
	s.keys = append(s.keys, key)
}

func (s *MultiSort) Len() int {
	return len(s.Tracks)
}

func (s *MultiSort) Swap(i, j int) {
	s.Tracks[i], s.Tracks[j] = s.Tracks[j], s.Tracks[i]
}

func (s *MultiSort) Less(i, j int) bool {
	x, y := s.Tracks[i], s.Tracks[j]
	for i := len(s.keys) - 1; i >= 0; i-- {
		if r := s.keys[i](x, y); r != 0 {
			return r < 0
		}
	}
	return false
}

func TitleCompare(x, y *Track) int {
	if x.Title < y.Title {
		return -1
	}
	if x.Title > y.Title {
		return 1
	}
	return 0
}

func ArtistCompare(x, y *Track) int {
	if x.Artist < y.Artist {
		return -1
	}
	if x.Artist > y.Artist {
		return 1
	}
	return 0
}

func AlbumCompare(x, y *Track) int {
	if x.Album < y.Album {
		return -1
	}
	if x.Album > y.Album {
		return 1
	}
	return 0
}

func YearCompare(x, y *Track) int {
	if x.Year < y.Year {
		return -1
	}
	if x.Year > y.Year {
		return 1
	}
	return 0
}

func LengthCompare(x, y *Track) int {
	if x.Length < y.Length {
		return -1
	}
	if x.Length > y.Length {
		return 1
	}
	return 0
}
