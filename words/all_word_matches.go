package words

import (
	"fmt"
	"sync"
)

type AllWordMatches struct {
	Words            *WordMap
	PackedMatch      *PackedMatchMatrix
	ExactPackedMatch map[uint32]uint8
}

type PackedMatchMatrix struct {
	rows []map[uint16]uint8
}

func NewPackedMatchMatrix() *PackedMatchMatrix {
	size := 256 * 256
	rows := make([]map[uint16]uint8, size)
	for i := range rows {
		rows[i] = make(map[uint16]uint8)
	}
	return &PackedMatchMatrix{rows: rows}
}

func (p *PackedMatchMatrix) Set(i, j uint16, value uint8) {
	// The indexes are very full in the bottom byte, and a bit sparse in the top byte,
	// so we extract the bottom byte of i and j to form the row index, and then use a map
	// for the other bytes.
	lsb := (i&0xFF)<<8 | (j & 0xFF)
	msb := (i & 0xFF00) | ((j >> 8) & 0xFF)
	p.rows[lsb][msb] = value
}

func (p *PackedMatchMatrix) Get(i, j uint16) uint8 {
	top := (i&0xFF)<<8 | (j & 0xFF)
	bottom := (i & 0xFF00) | ((j >> 8) & 0xFF)
	return p.rows[top][bottom]
}

func (p *PackedMatchMatrix) Len() int {
	count := 0
	for _, row := range p.rows {
		count += len(row)
	}
	return count
}

type allResult struct {
	i, j       uint16
	match      uint8
	match2     uint8
	exactMatch uint8
}

func GetAllWordMatches(wordMap *WordMap) *AllWordMatches {
	// (n) * (n-1)/2*2 because we track both (i,j) and (j,i), but we don't track (i,i)
	// Track all matches from a to b, and b to a, but ignore (i,i) since those are all exact matches
	matches := NewPackedMatchMatrix()
	// Exact matches are symmetric so we only store the (i,j) where i<j
	expectedLen := len(wordMap.Words) * (len(wordMap.Words) - 1) / 2
	exactMatches := make(map[uint32]uint8, expectedLen)
	words := wordMap.Words
	results := make(chan []allResult)
	done := make(chan struct{})
	wg := sync.WaitGroup{}

	go collectResults(results, matches, exactMatches, done)
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

	if len(exactMatches) != expectedLen {
		panic(fmt.Sprintf("expected %d matches, got %d", expectedLen, len(exactMatches)))
	}
	return &AllWordMatches{
		Words:            wordMap,
		PackedMatch:      matches,
		ExactPackedMatch: exactMatches,
	}
}

func collectResults(results chan []allResult, matches *PackedMatchMatrix, exactMatches map[uint32]uint8, done chan struct{}) {
	for res := range results {
		for _, r := range res {
			matches.Set(r.i, r.j, r.match)
			matches.Set(r.j, r.i, r.match2)
			key := (uint32(r.i) << 16) | uint32(r.j)
			exactMatches[key] = r.exactMatch
		}
	}
	close(done)
}

func processIthWords(i int, words []string) []allResult {
	results := make([]allResult, 0, len(words)-i-1)
	a := words[i]
	for j := i + 1; j < len(words); j++ {
		b := words[j]
		packed, exactPacked := GetPackedMatches(a, b)
		packed2, exactPacked2 := GetPackedMatches(b, a)
		if exactPacked2 != exactPacked {
			panic("exact packed match should be symmetric")
		}
		results = append(results, allResult{
			i:          uint16(i),
			j:          uint16(j),
			match:      packed,
			match2:     packed2,
			exactMatch: exactPacked,
		})
	}
	return results
}

func (awm *AllWordMatches) Match(i, j uint16) (uint8, uint8) {
	if i == j {
		return AllMatchedPacked, AllExactMatchedPacked
	}
	packed := awm.PackedMatch.Get(i, j)
	key := (uint32(i) << 16) | uint32(j)
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
				// missesAB, missesBA := FindMisses(a, b, match)
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
	var missesA uint8
	var missesB uint8
	return missesA, missesB
}
