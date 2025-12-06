package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jameinel/wordle-solver/words"
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

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
