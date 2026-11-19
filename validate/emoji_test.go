package validate

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"
)

// 絵文字の並びを検索表で探した結果が、emojiPattern を正規表現として当てた結果と同じであることを確かめる
func TestFindEmojiIndexesMatchesRegex(t *testing.T) {
	re := regexp.MustCompile(emojiPattern)

	var sequences []string
	for _, alternative := range strings.Split(emojiPattern, "|") {
		sequences = append(sequences, decodeEmojiAlternative(alternative))
	}

	// 絵文字の前後に置く文字。数字と # はキーキャップ、ZWJ と異体字セレクタと肌の色は絵文字の並びの途中に現れる文字である
	affixes := []string{"あ", "a", "1", "#", "‍", "️", "\U0001F3FD"}

	inputs := []string{
		strings.Join(sequences, ""),
		strings.Join(sequences, "‍"),
		strings.Join(sequences, " "),
	}
	for _, sequence := range sequences {
		inputs = append(inputs, sequence)
		_, lastSize := utf8.DecodeLastRuneInString(sequence)
		inputs = append(inputs, sequence[:len(sequence)-lastSize])
		for _, affix := range affixes {
			inputs = append(inputs, affix+sequence+affix)
		}
	}

	for _, input := range inputs {
		want := re.FindAllStringIndex(input, -1)
		got := findEmojiIndexes(input)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("input %+q: got %v, want %v", input, got, want)
		}
	}
}
