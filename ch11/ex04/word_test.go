package word

import (
	"math/rand"
	"testing"
	"time"
)

var NonLetters = []rune{' ', ',', '.', '　', '、', '。'}

func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25)
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000))
		runes[i] = r
		runes[n-1-i] = r
	}
	// 非文字をランダムに追加
	n = rng.Intn(25)
	for i := 0; i < n; i++ {
		idx := rng.Intn(len(runes) + 1)
		r := NonLetters[rng.Intn(len(NonLetters))]

		var tmp []rune
		tmp = append(tmp, runes[:idx]...)
		tmp = append(tmp, r)
		runes = append(tmp, runes[idx:]...)
	}
	return string(runes)
}

func TestRandomPalindromes(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomPalindrome(rng)
		if !IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = false", p)
		}
	}
}
