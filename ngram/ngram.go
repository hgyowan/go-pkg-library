package ngram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/hgyowan/go-pkg-library/envs"
	"github.com/samber/lo"
	"strings"
	"unicode"
)

func isKorean(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Hangul, r) {
			return true
		}
	}
	return false
}

func tokenizeNgram(input string, minGram, maxGram int) []string {
	runes := []rune(input)
	var tokens []string
	for length := minGram; length <= maxGram; length++ {
		for i := 0; i <= len(runes)-length; i++ {
			tokens = append(tokens, string(runes[i:i+length]))
		}
	}

	return lo.Uniq(tokens)
}

func hmacToken(token string) string {
	h := hmac.New(sha256.New, []byte(envs.SecretKey))
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateHmacTokens(input string) []string {
	input = strings.ToLower(input)
	var tokens []string
	if isKorean(input) {
		tokens = tokenizeNgram(input, 1, 3)
	} else {
		tokens = tokenizeNgram(input, 2, 4)
	}

	var encrypted []string
	for _, t := range tokens {
		encrypted = append(encrypted, hmacToken(t))
	}
	return encrypted
}
