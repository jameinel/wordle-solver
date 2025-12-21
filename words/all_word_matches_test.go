package words

import (
	"sort"
	"testing"
)

func TestAllWordMatches_SmallList(t *testing.T) {
	words := []string{"apple", "pleap", "aplep", "zzzzz", "abcde", "abedf"}
	sort.Strings(words)
	wordMap := NewWordMap(words)
	allMatches := GetAllWordMatches(wordMap)

	for i, a := range words {
		for j := i; j < len(words); j++ {
			b := words[j]
			expectedPacked, expectedExactPacked := GetPackedMatches(a, b)
			packed, exactPacked := allMatches.Match(uint16(i), uint16(j))
			if packed != expectedPacked || exactPacked != expectedExactPacked {
				t.Errorf("Mismatch for %s vs %s: got packed=%d, exactPacked=%x; expected packed=%d, exactPacked=%x",
					a, b, packed, exactPacked, expectedPacked, expectedExactPacked)
			}
			packed2, exactPacked2 := allMatches.Match(uint16(j), uint16(i))
			expectedPacked2, expectedExactPacked2 := GetPackedMatches(b, a)
			if packed2 != expectedPacked2 || exactPacked2 != expectedExactPacked2 {
				t.Errorf("Mismatch for %s vs %s: got packed=%d, exactPacked=%x; expected packed=%d, exactPacked=%x",
					b, a, packed2, exactPacked, expectedPacked2, expectedExactPacked2)
			}
			if expectedExactPacked != expectedExactPacked2 {
				t.Errorf("Asymmetry in exact packed matches for %s vs %s: %x vs %x",
					a, b, expectedExactPacked, expectedExactPacked2)
			}
		}
	}
}
