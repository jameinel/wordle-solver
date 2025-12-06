package main

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/jameinel/wordle-solver/words"
)

func generateSmallWords(seed int64, fraction float64) {
	inFile := "words/words.json"
	outFile := "words/small_words.json"

	allWords, err := words.ReadWordsFile(inFile)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", inFile, err)
		os.Exit(1)
	}
	r := rand.New(rand.NewSource(seed))
	outWords := &words.Words{
		Solutions:  grabFraction(r, allWords.Solutions, fraction),
		Dictionary: grabFraction(r, allWords.Dictionary, fraction),
	}
	out, err := json.MarshalIndent(outWords, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outFile, out, 0644); err != nil {
		fmt.Printf("Error writing %s: %v\n", outFile, err)
		os.Exit(1)
	}
	fmt.Printf("Generated %s\n", outFile)
}

func grabFraction(r *rand.Rand, arr []string, fraction float64) []string {
	n := int(math.Floor(float64(len(arr)) * fraction))
	if n < 1 {
		n = 1
	}
	r.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr[:n]
}

func cryptoSeed() int64 {
	var b [8]byte
	_, err := cryptorand.Read(b[:])
	if err != nil {
		fmt.Printf("Error generating random seed: %v\n", err)
		os.Exit(1)
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}
