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

var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import wcup data from a compressed archive",
	Long: `Import results and participant predictions from a gzip-compressed tar archive.
This replaces all current data.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !ReadConfirmation("¿Importar datos? Esto reemplazará todos los datos actuales") {
			return
		}

		src := args[0]
		if err := lib.ImportData(src); err != nil {
			fmt.Fprintf(os.Stderr, "failed to import data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Datos importados")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
