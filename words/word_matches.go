package words

type LetterMatch int8

const (
	MatchUnknown LetterMatch = 0
	MatchNone    LetterMatch = 1
	MatchNear    LetterMatch = 2
	MatchExact   LetterMatch = 3
)

const (
	AllMatchedPacked      uint8 = 243 - 1 // 3^5 - 1
	AllExactMatchedPacked uint8 = 0x1F    // 2^5 - 1
)

func (lm LetterMatch) String() string {
	switch lm {
	case MatchUnknown:
		return "Unknown"
	case MatchNone:
		return "None"
	case MatchNear:
		return "Near"
	case MatchExact:
		return "Exact"
	default:
		return "Invalid"
	}
}

func GetWordMatches(test, goal string) []LetterMatch {
	if len(test) != len(goal) {
		return nil
	}
	matches := make([]LetterMatch, len(test))
	// unused := make(map[uint8]int)
	var unused [26]uint8
	for i := 0; i < len(goal); i++ {
		c := goal[i]
		if c == test[i] {
			matches[i] = MatchExact
		} else {
			cc := c - 'a'
			unused[cc] = unused[cc] + 1
		}
	}
	for i := 0; i < len(test); i++ {
		c := test[i]
		if c != goal[i] {
			cc := c - 'a'
			if count := unused[cc]; count > 0 {
				matches[i] = MatchNear
				unused[cc] = count - 1
			} else {
				matches[i] = MatchNone
			}
		}
	}
	return matches
}

func GetPackedMatches(test, goal string) (uint8, uint8) {
	matches := GetWordMatches(test, goal)
	if matches == nil {
		return 0xFF, 0xFF
	}
	var packed uint8
	var exactPacked uint8
	for _, match := range matches {
		packed *= 3
		exactPacked *= 2
		if match == MatchUnknown {
			// Should not happen
			return 0xFE, 0x00
		}
		packed += uint8(match) - 1
		if match == MatchExact {
			exactPacked += 1
		}
	}
	return packed, exactPacked
}
