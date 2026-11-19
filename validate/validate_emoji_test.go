package validate

import (
	"strings"
	"testing"
)

// 期待値は、PHP の twitter-text（nojimage/twitter-text-php v3.4.0）の Parser::parseTweet で同じ文を数えた結果である
func TestParseTweetWeightedLengthWithEmoji(t *testing.T) {
	tests := map[string]struct {
		text string
		want int
	}{
		"結合した絵文字は2文字と数える":       {text: "👨‍👩‍👧", want: 2},
		"肌の色付きの絵文字は2文字と数える":     {text: "👍🏽", want: 2},
		"国旗は2文字と数える":            {text: "🇯🇵", want: 2},
		"キーキャップは2文字と数える":        {text: "1️⃣", want: 2},
		"1文字の絵文字は2文字と数える":       {text: "😀", want: 2},
		"絵文字として扱う記号は2文字と数える":    {text: "©", want: 2},
		"前後の英数字は1文字ずつ数える":       {text: "a👨‍👩‍👧b", want: 4},
		"絵文字とURLと日本語と英数字が混ざった文": {text: "絵文字👍🏽と https://example.com/path と英数字abc🇯🇵", want: 48},
		"絵文字の連なりだけで280文字":       {text: strings.Repeat("👨‍👩‍👧", 140), want: 280},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := ParseTweet(tt.text)
			if got.WeightedLength != tt.want {
				t.Errorf("got %d, want %d", got.WeightedLength, tt.want)
			}
		})
	}
}
