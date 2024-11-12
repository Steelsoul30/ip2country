package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create-db",
	Aliases: []string{"gdb", "cdb", "create-database", "generate-database"},
	Short:   "Create database",
	Long:    "Generate database to serialize the database schema to a file in preparation for runtime use.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Creating database...\n")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
