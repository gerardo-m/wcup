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
	Long:  `Interactive editor for a participant's World Cup predictions.`,
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
	predictionCmd.Flags().StringVarP(&predictionParticipant, "participant", "p", "", "Participant name")
}
