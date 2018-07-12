package echo

import (
  "testing"
)

var ARGS = []string{"foo", "bar", "baz", "hoge", "fuga", "piyo"}

func BenchmarkCatStrs(b *testing.B) {
  for i := 0; i < b.N; i++ {
    catStrs(ARGS)
  }
}

func BenchmarkCatStrsJoin(b *testing.B) {
  for i := 0; i < b.N; i++ {
    catStrsJoin(ARGS)
  }
}
