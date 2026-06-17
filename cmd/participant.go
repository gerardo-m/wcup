/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/gerardo-m/wcup/lib"
	"github.com/spf13/cobra"
)

// participantCmd represents the participant command
var participantCmd = &cobra.Command{
	Use:   "participant",
	Short: "List registered participants",
	Long:  `Display all participants with saved prediction files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(lib.Participants) == 0 {
			fmt.Println("No hay participantes registrados")
			return
		}

		for _, participant := range lib.Participants {
			fmt.Println(participant.Name)
		}
	},
}

func init() {
	showCmd.AddCommand(participantCmd)
}
