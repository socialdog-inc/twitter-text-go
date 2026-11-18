package validate

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// emojiSequencesByFirstRune maps a first rune to the sequences in emojiPattern that start with it.
// They keep the order of emojiPattern, because regexp alternation takes the leftmost alternative that matches.
var emojiSequencesByFirstRune = buildEmojiSequencesByFirstRune(emojiPattern)

func buildEmojiSequencesByFirstRune(pattern string) map[rune][]string {
	table := map[rune][]string{}
	for _, alternative := range strings.Split(pattern, "|") {
		sequence := decodeEmojiAlternative(alternative)
		first, _ := utf8.DecodeRuneInString(sequence)
		table[first] = append(table[first], sequence)
	}
	return table
}

// decodeEmojiAlternative converts an alternative such as `\x{1f468}\x{200d}` into the string it matches.
// It panics on any other syntax, so that a change in the pattern does not silently change the result.
func decodeEmojiAlternative(alternative string) string {
	var b strings.Builder
	for rest := alternative; rest != ""; {
		end := strings.IndexByte(rest, '}')
		if !strings.HasPrefix(rest, `\x{`) || end < 0 {
			panic(fmt.Sprintf("validate: unsupported syntax in emoji pattern: %q", alternative))
		}
		code, err := strconv.ParseUint(rest[len(`\x{`):end], 16, 32)
		if err != nil {
			panic(fmt.Sprintf("validate: invalid code point in emoji pattern: %q", alternative))
		}
		b.WriteRune(rune(code))
		rest = rest[end+1:]
	}
	if b.Len() == 0 {
		panic("validate: empty alternative in emoji pattern")
	}
	return b.String()
}

// findEmojiIndexes returns the byte ranges of the emoji sequences in text, scanning from the start
// without overlaps. The result is the same as regexp.MustCompile(emojiPattern).FindAllStringIndex(text, -1).
func findEmojiIndexes(text string) [][]int {
	var indexes [][]int
	for offset := 0; offset < len(text); {
		r, size := utf8.DecodeRuneInString(text[offset:])
		if end := matchEmojiSequence(text, offset, r); end > 0 {
			indexes = append(indexes, []int{offset, end})
			offset = end
			continue
		}
		offset += size
	}
	return indexes
}

// matchEmojiSequence returns the end offset of the emoji sequence that starts at offset, or 0 if there is none.
func matchEmojiSequence(text string, offset int, first rune) int {
	for _, sequence := range emojiSequencesByFirstRune[first] {
		if strings.HasPrefix(text[offset:], sequence) {
			return offset + len(sequence)
		}
	}
	return 0
}
