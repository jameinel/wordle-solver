package main

import (
	"fmt"
	"os"

	"github.com/jameinel/wordle-solver/words"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "wordle",
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
	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a small_words.json with random 10% of La and Ta (optionally seeded)",
		Run: func(cmd *cobra.Command, args []string) {
			if seed == 0 {
				seed = cryptoSeed()
			}
			generateSmallWords(seed)
		},
	}

	generateCmd.Flags().Int64Var(&seed, "seed", 0, "Seed for random number generator (default 0)")
	rootCmd.AddCommand(generateCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
