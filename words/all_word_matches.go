package words

import (
	"fmt"
	"sync"
)

type AllWordMatches struct {
	Words            *WordMap
	PackedMatch      map[uint32]uint8
	ExactPackedMatch map[uint32]uint8
}

type allResult struct {
	key        uint32
	key2       uint32
	match      uint8
	match2     uint8
	exactMatch uint8
}

func processIthWords(i int, words []string) []allResult {
	results := make([]allResult, 0, len(words)-i-1)
	a := words[i]
	for j := i + 1; j < len(words); j++ {
		// fmt.Printf("%d,%d\n", ii, j)
		b := words[j]
		key := (uint32(i) << 16) | uint32(j)
		packed, exactPacked := GetPackedMatches(a, b)
		key2 := (uint32(j) << 16) | uint32(i)
		packed2, exactPacked2 := GetPackedMatches(b, a)
		if exactPacked2 != exactPacked {
			panic("exact packed match should be symmetric")
		}
		results = append(results, allResult{
			key:        key,
			key2:       key2,
			match:      packed,
			match2:     packed2,
			exactMatch: exactPacked,
		})
	}
	return results
}

func GetAllWordMatches(wordMap *WordMap) *AllWordMatches {
	// (n) * (n-1)/2*2 because we track both (i,j) and (j,i), but we don't track (i,i)
	expectedLen := len(wordMap.Words) * (len(wordMap.Words) - 1)
	// Track all matches from a to b, and b to a, but ignore (i,i) since those are all exact matches
	matches := make(map[uint32]uint8, expectedLen)
	// Exact matches are symmetric so we only store the (i,j) where i<j
	exactMatches := make(map[uint32]uint8, len(wordMap.Words))
	words := wordMap.Words
	results := make(chan []allResult)
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	go func() {
		for res := range results {
			for _, r := range res {
				matches[r.key] = r.match
				matches[r.key2] = r.match2
				exactMatches[r.key] = r.exactMatch
			}
		}
		close(done)
	}()
	for i := 0; i < len(words); i++ {
		wg.Add(1)
		go func(ii int) {
			res := processIthWords(ii, words)
			results <- res
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(results)
	<-done

	if len(matches) != expectedLen {
		panic(fmt.Sprintf("expected %d matches got %d", expectedLen, len(matches)))
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

type AllSmallWordMatches struct {
	Words   *SmallWordMap
	Matches map[uint32]uint8
}

type smallResult struct {
	key   uint32
	match uint8
}

func GetAllSmallWordMatches(wordMap *SmallWordMap) *AllSmallWordMatches {
	// Exact matches are symmetric so we only store the (i,j) where i<j
	expectedLen := len(wordMap.Words) * (len(wordMap.Words) - 1) / 2
	matches := make(map[uint32]uint8, expectedLen)
	words := wordMap.SmallWords
	results := make(chan []smallResult)
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	go func() {
		for res := range results {
			for _, r := range res {
				matches[r.key] = r.match
			}
		}
		close(done)
	}()

	for i := 0; i < len(words); i++ {
		wg.Add(1)
		go func(ii int) {
			res := make([]smallResult, 0, len(words)-ii-1)
			a := words[ii]
			for j := ii + 1; j < len(words); j++ {
				b := words[j]
				match := MatchPackedWords(a, b)
				key := (uint32(i) << 16) | uint32(j)
				res = append(res, smallResult{
					key:   key,
					match: match,
				})
				missedA, missedB := FindMisses(a, b, match)
				_ = missedA
				_ = missedB
			}
			results <- res
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(results)
	<-done

	if len(matches) != expectedLen {
		panic("unexpected number of matches computed")
	}
	return &AllSmallWordMatches{
		Words:   wordMap,
		Matches: matches,
	}
}

func maskMatchedLetters(match uint8) uint32 {
	var mask uint32
	// Map bit 0 to bit 0-5, bit 1 to 6-11, bit 2 to 12-17, bit 3 to 18-23, bit 4 to 24-29
	mask = uint32(match)
	mask = 0 |
		((mask & 0x10) << 20) |
		((mask & 0x08) << 15) |
		((mask & 0x04) << 10) |
		((mask & 0x02) << 5) |
		((mask & 0x01) << 0)
	mask *= 0x3F
	return (^mask) & 0x3FFFFFFF
}

func FindMisses(packedA, packedB uint32, match uint8) (uint8, uint8) {
	mask := maskMatchedLetters(match)
	unmatchedA := packedA & mask
	unmatchedB := packedB & mask
	unusedA := make(map[uint8]int)
	unusedB := make(map[uint8]int)
	for i := 0; i < 5; i++ {
		c := uint8(unmatchedA>>(6*i)) & 0x3F
		if c != 0 {
			unusedA[c] = unusedA[c] + 1
		}
		c = uint8(unmatchedB>>(6*i)) & 0x3F
		if c != 0 {
			unusedB[c] = unusedB[c] + 1
		}
	}
	var missesA, missesB uint8
	for i := 4; i >= 0; i-- {
		missesA <<= 1
		c := uint8(unmatchedA>>(6*i)) & 0x3F
		if c != 0 {
			count := unusedB[c]
			if count > 0 {
				unusedB[c] = count - 1
				missesA |= 1
			}
		}
		missesB <<= 1
		c = uint8(unmatchedB>>(6*i)) & 0x3F
		if c != 0 {
			count := unusedA[c]
			if count > 0 {
				unusedA[c] = count - 1
				missesB |= 1
			}
		}
	}
	return missesA, missesB
}
