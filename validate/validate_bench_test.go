package validate

import (
	"strings"
	"testing"
)

// 文字数を数える処理（絵文字の判定を含む）が遅くなっていないかを、変更の前後で比べるためのベンチマーク
func BenchmarkParseTweet(b *testing.B) {
	benchmarks := map[string]string{
		"Japanese25000": strings.Repeat("あいうえおかきくけこ", 2500),
		"Digits280":     strings.Repeat("0123456789", 28),
		"EmojiZWJ140":   strings.Repeat("👨‍👩‍👧", 140),
		"MixedWithURL":  strings.Repeat("絵文字👍🏽と https://example.com/path と英数字abc🇯🇵 ", 20),
	}
	for name, text := range benchmarks {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = ParseTweet(text)
			}
		})
	}
}
