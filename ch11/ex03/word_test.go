package word

import (
	"math/rand"
	"testing"
	"time"
	"unicode"
)

func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25)
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000))
		runes[i] = r
		runes[n-1-i] = r
	}
	return string(runes)
}

func randomNonPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25)
	runes := make([]rune, n)
	palindromeFlag := true

	for i := 0; i < (n+1)/2; i++ {
		r1 := rune(rng.Intn(0x1000))
		for !unicode.IsLetter(r1) {
			r1 = rune(rng.Intn(0x1000))
		}
		r2 := rune(rng.Intn(0x1000))
		for !unicode.IsLetter(r2) {
			r2 = rune(rng.Intn(0x1000))
		}
		if i != n-1-i && r1 != r2 {
			palindromeFlag = false
		}
		runes[i] = r1
		runes[n-1-i] = r2
	}

	if palindromeFlag {
		// 回文になってしまった場合は再生成
		return randomNonPalindrome(rng)
	} else {
		// 非文字をランダムに追加
		n = rng.Intn(25)
		for i := 0; i < n; i++ {
			idx := rng.Intn(len(runes) + 1)
			r := rune(rng.Intn(0x1000))
			for unicode.IsLetter(r) {
				r = rune(rng.Intn(0x1000))
			}
			var tmp []rune
			tmp = append(tmp, runes[:idx]...)
			tmp = append(tmp, r)
			runes = append(tmp, runes[idx:]...)
		}
		return string(runes)
	}
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

func TestRandomNonPalindromes(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		np := randomNonPalindrome(rng)
		if IsPalindrome(np) {
			t.Errorf("IsPalindrome(%q) = true", np)
		}
	}
}
