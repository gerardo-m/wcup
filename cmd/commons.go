package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadConfirmation(question string) bool {
	fmt.Print(question + " (y/n): ")
	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read confirmation: %v\n", err)
		os.Exit(1)
	}

	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("Respuesta inválida. Operación cancelada")
		return false
	}
}
