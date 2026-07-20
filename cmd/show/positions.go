/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/gerardo-m/wcup/lib"
	"github.com/spf13/cobra"
)

// positionsCmd represents the positions command
var positionsCmd = &cobra.Command{
	Use:   "positions",
	Short: "Show World Cup group standings",
	Long: `Calculate and print group standings from registered match
results: points, wins, draws, losses, and goal difference for each team.`,
	Example: `  wcup show positions`,
	Run: func(cmd *cobra.Command, args []string) {
		results := lib.MatchResultsByID()

		for i, group := range lib.Groups {
			fmt.Printf("Grupo %s\n", group.Name)
			fmt.Printf("%-4s %-6s %-4s %-3s %-3s %-3s %-4s\n", "Pos", "Equipo", "Pts", "W", "D", "L", "GD")
			for _, standing := range lib.GroupStandings(group, results) {
				fmt.Printf("%-4d %-6s %-4d %-3d %-3d %-3d %-4d\n",
					standing.Position,
					standing.Team.Abbr,
					standing.Points,
					standing.Wins,
					standing.Draws,
					standing.Losses,
					standing.GoalDifference,
				)
			}
			if i < len(lib.Groups)-1 {
				fmt.Println()
			}
		}
	},
}

func init() {
	showCmd.AddCommand(positionsCmd)
}
