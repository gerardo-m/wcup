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

var nameFlag string
var allFlag bool

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset <all|results|participant>",
	Short: "Reset results or participants",
	Long: `Clear registered World Cup data.

  all          Reset everything (not implemented yet)
  results      Clear all match results and knockout classifications
  participant  Remove one participant (--name) or every participant (--all)

Destructive actions ask for confirmation before proceeding.`,
	Example: `  wcup reset results
  wcup reset participant --name alice
  wcup reset participant --all
  wcup reset all`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: wcup reset <all|results|participant>")
			return
		}
		to_reset := args[0]
		switch to_reset {
		case "all":
			fmt.Println("Not implemented yet")
		case "results":
			runResetResultsOrExit()
		case "participant":
			runResetParticipantOrExit()
		default:
			fmt.Println("Usage: wcup reset <all|results|participant>")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "Participant name to reset")
	resetCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Reset all participants")
}

func runResetResults() error {
	if err := lib.ResetResults(); err != nil {
		return err
	}
	fmt.Println("Resultados reiniciados")
	return nil
}

func runResetResultsOrExit() {
	if !ReadConfirmation("¿Reiniciar todos los resultados?") {
		return
	}
	if err := runResetResults(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to reset results: %v\n", err)
		os.Exit(1)
	}
}

func runResetParticipantOrExit() {
	if allFlag {
		if !ReadConfirmation("¿Reiniciar todos los participantes?") {
			return
		}
		if err := lib.ResetAllParticipants(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to reset participants: %v\n", err)
			os.Exit(1)
		}
	}
	if nameFlag != "" {
		if !ReadConfirmation(fmt.Sprintf("¿Reiniciar el participante %s?", nameFlag)) {
			return
		}
		if err := lib.ResetParticipant(nameFlag); err != nil {
			fmt.Fprintf(os.Stderr, "failed to reset participant: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Usage: wcup reset participant --name <name> | --all")
	}
}
