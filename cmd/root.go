/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	enter "github.com/gerardo-m/wcup/cmd/enter"
	show "github.com/gerardo-m/wcup/cmd/show"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wcup",
	Short: "Calculator for World Cup predictions",
	Long: `Register official World Cup results, enter participant predictions,
and compare them to see who is leading the pool.

Use enter to record results or predictions interactively, show to display
standings and points, and export/import to back up or share your data.`,
	Example: `  wcup teams
  wcup enter results
  wcup enter prediction -p alice
  wcup show points
  wcup export backup.tar.gz`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(enter.Command())
	rootCmd.AddCommand(show.Command())
}
