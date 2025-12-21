package words

type AllWordMatches struct {
	Words            *WordMap
	PackedMatch      map[uint32]uint8
	ExactPackedMatch map[uint32]uint8
}

func GetAllWordMatches(wordMap *WordMap) *AllWordMatches {
	expectedLen := len(wordMap.Words) * (len(wordMap.Words) - 1)
	// Track all matches from a to b, and b to a, but ignore (i,i) since those are all exact matches
	matches := make(map[uint32]uint8, expectedLen)
	// Exact matches are symmetric so we only store the (i,j) where i<j
	exactMatches := make(map[uint32]uint8, len(wordMap.Words))
	words := wordMap.Words
	for i := 0; i < len(words); i++ {
		a := words[i]
		for j := i + 1; j < len(words); j++ {
			b := words[j]
			packed, exactPacked := GetPackedMatches(a, b)
			key := (uint32(i) << 16) | uint32(j)
			// Combine the two packed matches into a single byte
			matches[key] = packed
			exactMatches[key] = exactPacked
			packed2, exactPacked2 := GetPackedMatches(b, a)
			if exactPacked2 != exactPacked {
				panic("exact packed match should be symmetric")
			}
			key2 := (uint32(j) << 16) | uint32(i)
			matches[key2] = packed2
		}
	}
	if len(matches) != expectedLen {
		panic("unexpected number of matches computed")
	}
	return &AllWordMatches{
		Words:            wordMap,
		PackedMatch:      matches,
		ExactPackedMatch: exactMatches,
	}
}

func (awm *AllWordMatches) Match(i, j uint16) (uint8, uint8) {
	if i == j {
		return AllMatchedPacked, AllExactMatchedPacked
	}
	key := (uint32(i) << 16) | uint32(j)
	packed := awm.PackedMatch[key]
	if i > j {
		key = (uint32(j) << 16) | uint32(i)
	}
	exactPacked := awm.ExactPackedMatch[key]
	return packed, exactPacked
}
