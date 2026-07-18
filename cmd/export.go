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

var exportCmd = &cobra.Command{
	Use:   "export <file>",
	Short: "Export wcup data to a compressed archive",
	Long:  `Export results and participant predictions to a gzip-compressed tar archive.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dest := args[0]
		if err := lib.ExportData(dest); err != nil {
			fmt.Fprintf(os.Stderr, "failed to export data: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Datos exportados a %s\n", dest)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
