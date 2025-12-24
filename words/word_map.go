package words

import (
	"sort"
)

type WordMap struct {
	Words   []string
	Offsets map[string]int16
}

func NewWordMap(words []string) *WordMap {
	// TODO: do we want to vet that the words are all 5 letters
	offsets := make(map[string]int16)
	for i, word := range words {
		offsets[word] = int16(i)
	}
	return &WordMap{
		Words:   words,
		Offsets: offsets,
	}
}

func AllWordsToWordMap(allWords *Words) *WordMap {
	combined := make([]string, len(allWords.Solutions)+len(allWords.Dictionary))
	copy(combined, allWords.Solutions)
	copy(combined[len(allWords.Solutions):], allWords.Dictionary)
	sort.Strings(combined)
	return NewWordMap(combined)
}
