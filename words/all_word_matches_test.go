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

func TestMaskMatchedLetters(t *testing.T) {
	tests := []struct {
		name    string
		matched uint8
		mask    uint32
	}{
		{"none", 0b00000, 0b00_111111_111111_111111_111111_111111},
		{"first", 0b10000, 0b00_000000_111111_111111_111111_111111},
		{"second", 0b01000, 0b00_111111_000000_111111_111111_111111},
		{"third", 0b00100, 0b00_111111_111111_000000_111111_111111},
		{"fourth", 0b00010, 0b00_111111_111111_111111_000000_111111},
		{"fifth", 0b00001, 0b00_111111_111111_111111_111111_000000},
		{"all", 0b11111, 0b00_000000_000000_000000_000000_000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			masked := maskMatchedLetters(tt.matched)
			if masked != tt.mask {
				t.Errorf("MaskMatchedLetters(%05b)\nwant 0b%032b\n got 0b%032b", tt.matched, tt.mask, masked)
			}
		})
	}
}

func TestAllWords(t *testing.T) {
	filename := "small_words.json"
	// filename := "words.json"
	all, err := ReadWordsFile(filename)
	if err != nil {
		t.Fatalf("Error reading words file %s: %v", all, err)
	}
	wordMap := AllWordsToWordMap(all)
	allMatches := GetAllWordMatches(wordMap)

	// Spot check some known matches
	tests := []struct {
		a, b   string
		packed uint8
		exact  uint8
	}{
		{"seize", "seize", 81*2 + 27*2 + 9*2 + 3*2 + 1*2, 0b11111},
		{"seize", "stiff", 81*2 + 27*0 + 9*2 + 3*0 + 1*0, 0b10100},
		{"seize", "films", 81*1 + 27*0 + 9*1 + 3*0 + 1*0, 0b00000},
	}
	for _, tt := range tests {
		i, ok := wordMap.Offsets[tt.a]
		if !ok {
			t.Fatalf("Word %q not found in word map", tt.a)
		}
		j, ok := wordMap.Offsets[tt.b]
		if !ok {
			t.Fatalf("Word %q not found in word map", tt.b)
		}
		packed, exact := allMatches.Match(i, j)
		if packed != tt.packed || exact != tt.exact {
			t.Errorf("Mismatch for %s vs %s: got packed=%d, exact=%05b; expected packed=%d, exact=%05b",
				tt.a, tt.b, packed, exact, tt.packed, tt.exact)
		}
	}
}

func TestFindMisses(t *testing.T) {
	tests := []struct {
		name             string
		a, b             string
		match            uint8
		missesA, missesB uint8
	}{
		{"identical", "apple", "apple", 0b11111, 0b00000, 0b00000},
		{"anagram", "apple", "pleap", 0b00000, 0b11111, 0b11111},
		{"repeated letter", "apple", "appla", 0b11110, 0b00000, 0b00000},
		{"all_diff", "apple", "zzzzz", 0b00000, 0b00000, 0b00000},
		{"swap", "apple", "apelp", 0b11010, 0b00101, 0b00101},
		{"one diff", "apple", "aeplf", 0b10110, 0b00001, 0b01000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pa := packWordToUint32(tt.a)
			pb := packWordToUint32(tt.b)
			match := MatchPackedWords(pa, pb)
			if match != tt.match {
				t.Errorf("MatchPackedWords(%s, %s): got %05b, want %05b", tt.a, tt.b, match, tt.match)
			}
			missA, missB := FindMisses(pa, pb, match)
			if missA != tt.missesA {
				t.Errorf("FindMisses(%s, %s): missA got %05b, want %05b", tt.a, tt.b, missA, tt.missesA)
			}
			if missB != tt.missesB {
				t.Errorf("FindMisses(%s, %s): missB got %05b, want %05b", tt.a, tt.b, missB, tt.missesB)
			}
		})
	}
}
