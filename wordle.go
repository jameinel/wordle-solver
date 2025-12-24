package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jameinel/wordle-solver/words"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "wordle-solver",
		Short: "Simple wordle solver.",
		Long:  "A wordle solver that attempts to extend the information theorem methods that pick any word, to one that is aware of hard mode",
	}

	var solveCmd = &cobra.Command{
		Use:   "solve",
		Short: "Solve a wordle puzzle",
		Run: func(cmd *cobra.Command, args []string) {
			filename := viper.GetString("file")
			if filename == "" {
				fmt.Println("No filename provided. Use --file or -f to specify the input file.")
				os.Exit(1)
			}
			allWords, err := words.ReadWordsFile(filename)
			if err != nil {
				fmt.Printf("Error reading words: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Dictionary: %v\n", allWords.Dictionary)
			fmt.Printf("Solutions: %v\n", allWords.Solutions)
		},
	}

	solveCmd.Flags().StringP("file", "f", "", "Input JSON file")
	_ = viper.BindPFlag("file", solveCmd.Flags().Lookup("file"))
	rootCmd.AddCommand(solveCmd)

	var seed int64
	var fraction float64
	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a small_words.json with random 10% of La and Ta (optionally seeded)",
		Run: func(cmd *cobra.Command, args []string) {
			if seed == 0 {
				seed = cryptoSeed()
			}
			generateSmallWords(seed, fraction)
		},
	}

	generateCmd.Flags().Int64Var(&seed, "seed", 0, "Seed for random number generator (default 0)")
	generateCmd.Flags().Float64Var(&fraction, "fraction", 0.1, "Fraction of words to include (default 0.1)")
	rootCmd.AddCommand(generateCmd)

	var matchesCmd = &cobra.Command{
		Use:   "matches",
		Short: "Compute packed matches between all words",
		Run: func(cmd *cobra.Command, args []string) {
			filename := viper.GetString("file")
			if filename == "" {
				fmt.Println("No filename provided. Use --file or -f to specify the input file.")
				os.Exit(1)
			}
			allWords, err := words.ReadWordsFile(filename)
			if err != nil {
				fmt.Printf("Error reading words: %v\n", err)
				os.Exit(1)
			}

			wordMap := words.AllWordsToWordMap(allWords)

			tStart := time.Now()
			allMatches := words.GetAllWordMatches(wordMap)

			fmt.Printf("Computed %d matches in %s:\n", len(allMatches.ExactPackedMatch), time.Since(tStart))
		},
	}

	matchesCmd.Flags().StringP("file", "f", "", "Input JSON file")
	_ = viper.BindPFlag("file", matchesCmd.Flags().Lookup("file"))
	rootCmd.AddCommand(matchesCmd)

	var exactMatchesCmd = &cobra.Command{
		Use:   "exact-matches",
		Short: "Compute exact matches between all words",
		Run: func(cmd *cobra.Command, args []string) {
			filename := viper.GetString("file")
			if filename == "" {
				fmt.Println("No filename provided. Use --file or -f to specify the input file.")
				os.Exit(1)
			}
			allWords, err := words.ReadWordsFile(filename)
			if err != nil {
				fmt.Printf("Error reading words: %v\n", err)
				os.Exit(1)
			}

			wordMap := words.AllWordsToWordMap(allWords)
			smallWordMap := words.AllWordsToSmallWordMap(allWords)

			tStart := time.Now()
			smallWordMatches := words.GetAllSmallWordMatches(smallWordMap)
			tSmallMatches := time.Now()
			fmt.Printf("Computed %d small matches in %s\n", len(smallWordMatches.Matches), tSmallMatches.Sub(tStart))
			allMatches := words.GetAllWordMatches(wordMap)
			fmt.Printf("Computed %d all matches and %d exact matches in %s\n", allMatches.PackedMatch.Len(), len(allMatches.ExactPackedMatch), time.Since(tSmallMatches))
			if len(allMatches.ExactPackedMatch) != len(smallWordMatches.Matches) {
				fmt.Printf("Mismatch in number of exact matches: all=%d vs small=%d\n", len(allMatches.ExactPackedMatch), len(smallWordMatches.Matches))
				os.Exit(1)
			}
			for key, match := range smallWordMatches.Matches {
				allMatch := allMatches.ExactPackedMatch[key]
				if match != allMatch {
					i := uint16((key >> 16) & 0xFFFF)
					j := uint16(key & 0xFFFF)
					a := smallWordMap.Words[i]
					b := smallWordMap.Words[j]
					fmt.Printf("Mismatch in exact match for %q vs %q: all=%d vs small=%d\n", a, b, allMatch, match)
					os.Exit(1)
				}
			}
			fmt.Printf("All exact matches verified between small and all word maps.\n")
		},
	}

	exactMatchesCmd.Flags().StringP("file", "f", "", "Input JSON file")
	_ = viper.BindPFlag("file", exactMatchesCmd.Flags().Lookup("file"))
	rootCmd.AddCommand(exactMatchesCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
