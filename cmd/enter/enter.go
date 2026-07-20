/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// enterCmd represents the enter command
var enterCmd = &cobra.Command{
	Use:   "enter",
	Short: "Enter results or predictions interactively",
	Long: `Open an interactive editor to record World Cup match results
or a participant's predictions.

Use enter results for official outcomes, or enter prediction with a
participant name to fill in their forecast.`,
	Example: `  wcup enter results
  wcup enter prediction -p alice`,
}

func Command() *cobra.Command {
	return enterCmd
}
