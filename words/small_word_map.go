package words

import (
	"sort"
	"strings"
)

type SmallWordMap struct {
	Words      []string
	SmallWords []uint32
	Offsets    map[string]int16
}

func NewSmallWordMap(words []string) *SmallWordMap {
	offsets := make(map[string]int16)
	smallWords := make([]uint32, len(words))
	for i, word := range words {
		offsets[word] = int16(i)
		smallWords[i] = packWordToUint32(word)
	}
	return &SmallWordMap{
		Words:      words,
		SmallWords: smallWords,
		Offsets:    offsets,
	}
}

func packWordToUint32(word string) uint32 {
	var packed uint32
	if len(word) != 5 {
		panic("word must be 5 letters")
	}
	word = strings.ToLower(word)
	for i := 0; i < len(word); i++ {
		packed <<= 6
		packed |= uint32(word[i] - 'a' + 1)
	}
	return packed
}

func unpackWordFromUint32(packed uint32) string {
	var chars [5]byte
	for i := 4; i >= 0; i-- {
		chars[i] = byte((packed & 0x3F) + 'a' - 1)
		packed >>= 6
	}
	return string(chars[:])
}

type SmallWordMatches struct {
	Words   *SmallWordMap
	Matches map[uint32]uint8
}

// With 6 bits per letter, we can pack 5 letters into a 30-bit uint32
// This is the bottom 3 bits of each 6-bit sequence
const bot_three_bits uint32 = 0b000111_000111_000111_000111_000111

// The high bit of each 6-bit sequence
const high_bits uint32 = 0b100000_100000_100000_100000_100000

// After finding the matches in each sequence of 6-bits, we can use this multiplier
// to move the bits into the top 5 bits of a word, and then shift them back down
// The bits are at 29, 23, 17, 11, 1, moving them to 31, 30, 29, 28, 27 requires
// shifting by 2, 7, 12, 17, 22, respectively
const multiplier = 1<<2 | 1<<7 | 1<<12 | 1<<17 | 1<<22

// Find the exact matches from a to b using packed uint32 words
func MatchPackedWords(a, b uint32) uint8 {
	var a_xor_b uint32 = (a ^ b)
	// Now we want to find if any of the 6-bit sequences are zero
	tmp := a_xor_b | ((bot_three_bits & a_xor_b) << 3)
	tmp &= ^bot_three_bits
	tmp = tmp | (tmp << 2) | (tmp << 1)
	tmp &= high_bits
	// If any of the bits were set, then they will be set in the high bit
	// now we invert the bits to get 1 for matches
	tmp = high_bits - tmp
	// Now we need to move the bits down into the low 5 bits
	return uint8(((tmp * multiplier) >> 27) & 0x1F)
}

func AllWordsToSmallWordMap(allWords *Words) *SmallWordMap {
	combined := make([]string, len(allWords.Solutions)+len(allWords.Dictionary))
	copy(combined, allWords.Solutions)
	copy(combined[len(allWords.Solutions):], allWords.Dictionary)
	sort.Strings(combined)
	return NewSmallWordMap(combined)
}
