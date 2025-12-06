package words

import (
	"strings"
	"testing"
)

func TestReadFromJSON(t *testing.T) {
	jsonData := `{
		"Ta": ["apple", "banana"],
		"La": ["solution1", "solution2"]
	}`

	reader := strings.NewReader(jsonData)
	words, err := readFromJSON(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDictionary := []string{"apple", "banana"}
	expectedSolutions := []string{"solution1", "solution2"}

	if len(words.Dictionary) != len(expectedDictionary) {
		t.Errorf("expected Dictionary length %d, got %d", len(expectedDictionary), len(words.Dictionary))
	}
	for i, v := range expectedDictionary {
		if words.Dictionary[i] != v {
			t.Errorf("expected Dictionary[%d] = %q, got %q", i, v, words.Dictionary[i])
		}
	}

	if len(words.Solutions) != len(expectedSolutions) {
		t.Errorf("expected Solutions length %d, got %d", len(expectedSolutions), len(words.Solutions))
	}
	for i, v := range expectedSolutions {
		if words.Solutions[i] != v {
			t.Errorf("expected Solutions[%d] = %q, got %q", i, v, words.Solutions[i])
		}
	}
}
