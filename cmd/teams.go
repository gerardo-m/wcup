/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/gerardo-m/wcup/lib"
	"github.com/spf13/cobra"
)

// teamsCmd represents the teams command
var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "List World Cup groups and their teams",
	Long: `Display all group stage groups for the 2026 FIFA World Cup
with each team's abbreviation and full name.`,
	Example: `  wcup teams`,
	Run: func(cmd *cobra.Command, args []string) {
		for i, group := range lib.Groups {
			fmt.Printf("Group %s\n", group.Name)
			for _, team := range group.Teams {
				fmt.Printf("  %s  %s\n", team.Abbr, team.Name)
			}
			if i < len(lib.Groups)-1 {
				fmt.Println()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(teamsCmd)
}
