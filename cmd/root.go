/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/root-man/quiz/dummy"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Runs the quiz",
	Long:  `Runs the quiz`,
	Run: func(cmd *cobra.Command, args []string) {
		filepath, err := cmd.Flags().GetString("file")
		if err != nil {
			cobra.CompErrorln(fmt.Sprintf("filepath flag error: %s", err.Error()))
		}

		timer, err := cmd.Flags().GetInt("timer")
		if err != nil {
			cobra.CompErrorln(fmt.Sprintf("timer flag error: %s", err.Error()))
		}

		shuffle, err := cmd.Flags().GetBool("shuffle")
		if err != nil {
			cobra.CompErrorln(fmt.Sprintf("shuffle flag error: %s", err.Error()))
		}

		timeout := time.Duration(timer) * time.Second

		dummy.RunQuiz(dummy.NewOpts(dummy.WithProblemsFile(filepath), dummy.WithShuffle(shuffle), dummy.WithTimer(timeout)))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("file", "f", "problems.csv", "Specify path to the problems CSV file")
	rootCmd.Flags().IntP("timer", "t", 30, "Specify the time in seconds to run the quiz")
	rootCmd.Flags().BoolP("shuffle", "s", false, "Set the flag to shuffle the questions randomly")
}
