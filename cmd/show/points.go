/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/gerardo-m/wcup/lib"
	"github.com/spf13/cobra"
)

// pointsCmd represents the points command
var pointsCmd = &cobra.Command{
	Use:   "points",
	Short: "Show participant points",
	Long: `Calculate and display points for every participant by comparing
their predictions against the current official results.

The table includes totals and a breakdown by stage (matches, knockout
rounds, top scorer, and champion).`,
	Example: `  wcup show points`,
	Run: func(cmd *cobra.Command, args []string) {
		scores := lib.CalculateAllParticipantPoints()
		if len(scores) == 0 {
			fmt.Println("No hay participantes registrados")
			return
		}

		fmt.Printf("%-16s %-6s %-5s %-5s %-5s %-5s %-5s %-5s %-5s %-5s %-5s\n",
			"Participante", "Total", "Part.", "R32", "R16", "CF", "SF", "3er", "Fin", "Gol", "Cam")
		for _, score := range scores {
			fmt.Printf("%-16s %-6d %-5d %-5d %-5d %-5d %-5d %-5d %-5d %-5d %-5d\n",
				score.Name,
				score.Total,
				score.MatchPoints,
				score.Round32,
				score.Round16,
				score.QuarterFinal,
				score.SemiFinal,
				score.ThirdPlace,
				score.Final,
				score.TopScorer,
				score.Champion,
			)
		}
	},
}

func init() {
	showCmd.AddCommand(pointsCmd)
}
