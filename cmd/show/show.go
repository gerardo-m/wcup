/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display results, standings, or points",
	Long: `Display registered World Cup data: match results and knockout
classifications, group standings derived from those results, or the
points table for all participants.`,
	Example: `  wcup show results
  wcup show positions
  wcup show points`,
}

func Command() *cobra.Command {
	return showCmd
}
