/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/gerardo-m/wcup/lib"
	"github.com/spf13/cobra"
)

// participantCmd represents the participant command
var participantCmd = &cobra.Command{
	Use:   "participant",
	Short: "List registered participants",
	Long:  `Display all participants with saved prediction files.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 && args[0] == "add" {
			if len(args) != 2 {
				fmt.Println("Usage: wcup participant add <name>")
				return
			}
			name := args[1]
			if err := lib.EnsureParticipantFile(name); err != nil {
				fmt.Fprintf(os.Stderr, "failed to create participant file: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Participant added successfully")
			return
		}
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
	rootCmd.AddCommand(participantCmd)
}
