/*
Copyright © 2026 Gerardo Miranda <contact@gerardomiranda.dev>

*/
package main

import (
	"fmt"
	"os"

	"github.com/gerardo-m/wcup/cmd"
	"github.com/gerardo-m/wcup/lib"
)

func main() {
	lib.BuildSchedule()

	if err := lib.EnsureResultsFile(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create results file: %v\n", err)
		os.Exit(1)
	}
	if err := lib.LoadResults(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load results: %v\n", err)
		os.Exit(1)
	}
	if err := lib.EnsureParticipantsDir(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create participants directory: %v\n", err)
		os.Exit(1)
	}
	if err := lib.LoadParticipants(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load participants: %v\n", err)
		os.Exit(1)
	}

	cmd.Execute()
}
