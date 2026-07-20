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
	Short: "Show World Cup results",
	Long: `Print the currently registered World Cup results: group stage
match scores, knockout classifications, podium, and top scorer.`,
	Example: `  wcup show results`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Results:")
		fmt.Println("Match Results:")
		for _, matchResult := range lib.MatchResults {
			fmt.Printf("%s %d - %d %s\n", matchResult.Match.Team1.Name, matchResult.Team1Score, matchResult.Team2Score, matchResult.Match.Team2.Name)
		}
		fmt.Println("Dieciseisavos de final:")
		for _, team := range lib.RoundOf32 {
			fmt.Println(team.Name)
		}
		fmt.Println("Octavos de final:")
		for _, team := range lib.RoundOf16 {
			fmt.Println(team.Name)
		}
		fmt.Println("Cuartos de final:")
		for _, team := range lib.RoundOf8 {
			fmt.Println(team.Name)
		}
		fmt.Println("Semifinales:")
		for _, team := range lib.RoundOf4 {
			fmt.Println(team.Name)
		}
		fmt.Println("Final:")
		for _, team := range lib.RoundOf2 {
			fmt.Println(team.Name)
		}
		fmt.Println("Podio:")
		for _, team := range lib.Podium {
			fmt.Println(team.Name)
		}
		fmt.Println("Goleador:")
		fmt.Println(lib.TopScorer)
	},
}

func init() {
	showCmd.AddCommand(resultsCmd)
}
