/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// resultsCmd represents the results command
var resultsCmd = &cobra.Command{
	Use:   "results",
	Short: "Enter World Cup results interactively",
	Long:  `Interactive editor for match results and knockout classifications.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runResultsEditor(); err != nil {
			fmt.Fprintf(os.Stderr, "results editor failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	enterCmd.AddCommand(resultsCmd)
}
