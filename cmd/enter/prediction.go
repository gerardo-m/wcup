/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var predictionParticipant string

// predictionCmd represents the prediction command
var predictionCmd = &cobra.Command{
	Use:   "prediction",
	Short: "Enter participant predictions interactively",
	Long: `Open an interactive editor for a participant's World Cup
predictions: group stage scores, knockout picks, champion, and top scorer.

The participant must already exist (see wcup participant add). Pass their
name with --participant / -p.`,
	Example: `  wcup enter prediction -p alice
  wcup enter prediction --participant bob`,
	Run: func(cmd *cobra.Command, args []string) {
		if predictionParticipant == "" {
			fmt.Fprintln(os.Stderr, "usage: wcup enter prediction -p <participant>")
			os.Exit(1)
		}
		if err := runPredictionEditor(predictionParticipant); err != nil {
			fmt.Fprintf(os.Stderr, "prediction editor failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	enterCmd.AddCommand(predictionCmd)
	predictionCmd.Flags().StringVarP(&predictionParticipant, "participant", "p", "", "Participant whose predictions to edit")
}
