/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/gerardo-m/wcup/lib"
	"github.com/spf13/cobra"
)

// resultsCmd represents the results command
var resultsCmd = &cobra.Command{
	Use:   "results",
	Short: "Show saved match results and knockout predictions",
	Long:  `Display group stage match results and classified teams for each knockout phase.`,
	Run: func(cmd *cobra.Command, args []string) {
		resultsByMatchID := make(map[int]lib.MatchResult, len(lib.MatchResults))
		for _, result := range lib.MatchResults {
			resultsByMatchID[result.Match.Id] = result
		}

		for i, group := range lib.Groups {
			fmt.Printf("Grupo %s\n", group.Name)
			for _, match := range lib.MatchesByGroup[group.Name] {
				if result, ok := resultsByMatchID[match.Id]; ok {
					fmt.Printf("  %s %d - %d %s\n", match.Team1.Abbr, result.Team1Score, result.Team2Score, match.Team2.Abbr)
				} else {
					fmt.Printf("  %s vs %s  PENDIENTE\n", match.Team1.Abbr, match.Team2.Abbr)
				}
			}
			if i < len(lib.Groups)-1 {
				fmt.Println()
			}
		}

		fmt.Println()
		printTeamSection("RONDA DE 32", lib.RoundOf32)
		printTeamSection("OCTAVOS", lib.RoundOf16)
		printTeamSection("CUARTOS", lib.RoundOf8)
		printTeamSection("SEMIFINALES", lib.RoundOf4)
		printTeamSection("FINAL", lib.RoundOf2)
		printTeamSection("PODIO", lib.Podium)
	},
}

func printTeamSection(heading string, teams []lib.Team) {
	fmt.Println(heading)
	if len(teams) == 0 {
		fmt.Println("  PENDIENTE")
		return
	}
	for _, team := range teams {
		fmt.Printf("  %s  %s\n", team.Abbr, team.Name)
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(resultsCmd)
}
