package words

import (
	"fmt"
	"strings"
	"testing"
)

func TestPackUnpackRoundTrip(t *testing.T) {
	tests := []struct {
		word   string
		packed uint32
	}{
		{"aaaaa", 0b000001_000001_000001_000001_000001},
		{"hello", 0b001000_000101_001100_001100_001111},
		{"gamer", 0b000111_000001_001101_000101_010010},
		{"speed", 0b010011_010000_000101_000101_000100},
		{"think", 0b010100_001000_001001_001110_001011},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			packed := packWordToUint32(tt.word)
			unpacked := unpackWordFromUint32(packed)
			if unpacked != tt.word {
				t.Errorf("Round-trip for %q failed: got %q", tt.word, unpacked)
			}
			if packed != tt.packed {
				t.Errorf("Packing %q: expected 0b%030b, got 0b%030b", tt.word, tt.packed, packed)
			}
		})
	}
}

func formatBinarySixesWithUnderscores(n uint32) string {
	// Convert the integer to a binary string
	binStr := fmt.Sprintf("%030b", n)

	var parts []string
	// Iterate from the end of the string, taking 4 digits at a time
	for i := len(binStr); i > 0; i -= 6 {
		start := i - 6
		if start < 0 {
			start = 0
		}
		parts = append([]string{binStr[start:i]}, parts...)
	}

	// Join the parts with underscores
	return strings.Join(parts, "_")
}

func TestMatchPackedWords(t *testing.T) {
	tests := []struct {
		a, b    string
		matches uint8
	}{
		{"apple", "apple", 0b11111},
		{"abcde", "abxyz", 0b11000},
		{"abcde", "xyzde", 0b00011},
		{"axxxx", "ayyyy", 0b10000},
		{"baxxx", "cayyy", 0b01000},
		{"xxaxx", "yyayy", 0b00100},
		{"xxxax", "yyyay", 0b00010},
		{"xxxxa", "yyyya", 0b00001},
	}
	b2s := formatBinarySixesWithUnderscores
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s vs %s", tt.a, tt.b), func(t *testing.T) {
			packedA := packWordToUint32(tt.a)
			packedB := packWordToUint32(tt.b)
			narrow := MatchPackedWords(packedA, packedB)
			if narrow != tt.matches { // || narrow != tt.matches
				t.Errorf("Matching %q vs %q: want: 0b%05b\n got: 0b%05b\n  pA: 0b%s\n  pB: 0b%s",
					tt.a, tt.b, tt.matches, narrow, b2s(packedA), b2s(packedB))
			}
		})
	}
}

func TestNewSmallWordMap(t *testing.T) {
	words := []string{"apple", "aplep", "pleap", "zzzzz"}
	ints := []uint32{
		1<<24 | 16<<18 | 16<<12 | 12<<6 | 5,   // apple
		1<<24 | 16<<18 | 12<<12 | 5<<6 | 16,   // aplep
		16<<24 | 12<<18 | 5<<12 | 1<<6 | 16,   // pleap
		26<<24 | 26<<18 | 26<<12 | 26<<6 | 26, // zzzzz
	}
	swm := NewSmallWordMap(words)

	if len(swm.Words) != len(words) {
		t.Fatalf("expected %d words, got %d", len(words), len(swm.Words))
	}
	if len(swm.SmallWords) != len(words) {
		t.Fatalf("expected %d words, got %d", len(words), len(swm.SmallWords))
	}

	for i, word := range words {
		// Check Words slice contains the word at correct position
		if swm.Words[i] != word {
			t.Errorf("expected word at position %d: %q, got: %q", i, word, swm.Words[i])
		}
		if swm.SmallWords[i] != ints[i] {
			t.Errorf("word %s at position %d\nwant: 0b%030b\n got: 0b%030b", word, i, ints[i], swm.SmallWords[i])
		}

		// Check Offsets map contains correct offset
		offset, ok := swm.Offsets[word]
		if !ok {
			t.Errorf("word %q not found in Offsets map", word)
			continue
		}
		if int(offset) != i {
			t.Errorf("expected offset for %q: %d, got: %d", word, i, offset)
		}
	}
}
