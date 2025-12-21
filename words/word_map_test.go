package words

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewWordMapOffsets(t *testing.T) {
	words := []string{"apple", "banana", "cherry", "date"}
	wm := NewWordMap(words)

	// Verify Words slice matches input
	if len(wm.Words) != len(words) {
		t.Fatalf("expected %d words, got %d", len(words), len(wm.Words))
	}

	for i, word := range words {
		// Check Words slice contains the word at correct position
		if wm.Words[i] != word {
			t.Errorf("expected word at position %d: %q, got: %q", i, word, wm.Words[i])
		}

		// Check Offsets map contains correct offset
		offset, ok := wm.Offsets[word]
		if !ok {
			t.Errorf("word %q not found in Offsets map", word)
			continue
		}
		if int(offset) != i {
			t.Errorf("expected offset for %q: %d, got: %d", word, i, offset)
		}
	}
}

func TestAllWordsToWordMap_SortsAndCombines(t *testing.T) {
	allWords := &Words{
		Solutions:  []string{"pear", "apple"},
		Dictionary: []string{"banana", "date"},
	}
	expected := []string{"apple", "banana", "date", "pear"}
	wm := AllWordsToWordMap(allWords)

	if !reflect.DeepEqual(wm.Words, expected) {
		t.Errorf("expected sorted words %v, got %v", expected, wm.Words)
	}

	for i, word := range expected {
		offset, ok := wm.Offsets[word]
		if !ok {
			t.Errorf("word %q not found in Offsets", word)
			continue
		}
		if int(offset) != i {
			t.Errorf("expected offset for %q: %d, got: %d", word, i, offset)
		}
	}

	// Ensure the words are sorted
	if !sort.StringsAreSorted(wm.Words) {
		t.Errorf("words are not sorted: %v", wm.Words)
	}
}
