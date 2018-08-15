package uniq

func uniq(strs []string) []string {
  if len(strs) <= 0 {
    return strs
  }
  dups := 0
  prev := strs[0]
  for i := 1; i < len(strs); i++ {
    s := strs[i]
    if s == prev {
      dups++
    } else {
      strs[i-dups] = s
    }
    prev = s
  }
  return strs[:len(strs)-dups]
}
