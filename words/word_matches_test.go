package words

import (
	"testing"
)

func TestGetWordMatches(t *testing.T) {
	tests := []struct {
		a, b     string
		expected []LetterMatch
	}{
		// All exact matches
		{"apple", "apple", []LetterMatch{MatchExact, MatchExact, MatchExact, MatchExact, MatchExact}},
		// No matches
		{"apple", "zzzzz", []LetterMatch{MatchNone, MatchNone, MatchNone, MatchNone, MatchNone}},
		// Some near matches
		{"apple", "pleap", []LetterMatch{MatchNear, MatchNear, MatchNear, MatchNear, MatchNear}},
		// Mixed exact, near, none
		{"apple", "aplep", []LetterMatch{MatchExact, MatchExact, MatchNear, MatchNear, MatchNear}},
		// Repeated letters
		{"aabbc", "abbab", []LetterMatch{MatchExact, MatchNear, MatchExact, MatchNear, MatchNone}},
		{"abbab", "aabbc", []LetterMatch{MatchExact, MatchNear, MatchExact, MatchNear, MatchNone}},
		// Different lengths
		{"apple", "app", nil},
	}

	for _, tt := range tests {
		result := GetWordMatches(tt.a, tt.b)
		if len(result) != len(tt.expected) {
			if tt.expected == nil && result == nil {
				continue
			}
			t.Errorf("GetWordMatches(%q, %q): expected length %d, got %d", tt.a, tt.b, len(tt.expected), len(result))
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("GetWordMatches(%q, %q)[%d]: %q vs %q expected %s, got %s",
					tt.a, tt.b, i, tt.a[i], tt.b[i], tt.expected[i].String(), result[i].String())
			}
		}
	}
}

func TestGetPackedMatches(t *testing.T) {
	tests := []struct {
		a, b        string
		expected    []LetterMatch
		packed      uint8
		exactPacked uint8
	}{
		// The matches come out as 3-base number, so we can calculate the expected packed value
		// First match only
		{"azzzz", "ayyyy", []LetterMatch{MatchExact, MatchNone, MatchNone, MatchNone, MatchNone}, 81 * 2, 0x10},
		// Last match only
		{"zzzza", "yyyya", []LetterMatch{MatchNone, MatchNone, MatchNone, MatchNone, MatchExact}, 2, 0x01},
		// All exact matches
		{"apple", "apple", []LetterMatch{MatchExact, MatchExact, MatchExact, MatchExact, MatchExact}, 243 - 1, 0x1F},
		// No matches
		{"apple", "zzzzz", []LetterMatch{MatchNone, MatchNone, MatchNone, MatchNone, MatchNone}, 0, 0x00},
		// Some near matches
		{"apple", "pleap", []LetterMatch{MatchNear, MatchNear, MatchNear, MatchNear, MatchNear}, 81*1 + 27*1 + 9*1 + 3*1 + 1*1, 0x00},
		// Mixed exact, near, none
		{"apple", "aplep", []LetterMatch{MatchExact, MatchExact, MatchNear, MatchNear, MatchNear}, 81*2 + 27*2 + 9*1 + 3*1 + 1*1, 0x18},
		// Repeated letters
		{"aabbc", "abbab", []LetterMatch{MatchExact, MatchNear, MatchExact, MatchNear, MatchNone}, 81*2 + 27*1 + 9*2 + 3*1 + 1*0, 0x14},
		{"abbab", "aabbc", []LetterMatch{MatchExact, MatchNear, MatchExact, MatchNear, MatchNone}, 81*2 + 27*1 + 9*2 + 3*1 + 1*0, 0x14},
		// Non-symmetric matches
		{"abcde", "abedf", []LetterMatch{MatchExact, MatchExact, MatchNone, MatchExact, MatchNear}, 81*2 + 27*2 + 9*0 + 3*2 + 1*1, 0x1A},
		{"abedf", "abcde", []LetterMatch{MatchExact, MatchExact, MatchNear, MatchExact, MatchNone}, 81*2 + 27*2 + 9*1 + 3*2 + 1*0, 0x1A},
		// Different lengths
		{"apple", "app", nil, 0xFF, 0xFF},
	}

	for _, tt := range tests {
		result := GetWordMatches(tt.a, tt.b)
		if len(result) != len(tt.expected) {
			if tt.expected == nil && result == nil {
				continue
			}
			t.Errorf("GetWordMatches(%q, %q): expected length %d, got %d", tt.a, tt.b, len(tt.expected), len(result))
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("GetWordMatches(%q, %q)[%d]: %q vs %q expected %s, got %s",
					tt.a, tt.b, i, tt.a[i], tt.b[i], tt.expected[i].String(), result[i].String())
			}
		}
		packedResult, exactPackedResult := GetPackedMatches(tt.a, tt.b)
		if packedResult != tt.packed {
			t.Errorf("GetWordMatches(%q, %q): expected %d, got %d",
				tt.a, tt.b, tt.packed, packedResult)
		}
		if exactPackedResult != tt.exactPacked {
			t.Errorf("GetWordMatches(%q, %q): expected %d, got %d",
				tt.a, tt.b, tt.exactPacked, exactPackedResult)
		}
	}
}
