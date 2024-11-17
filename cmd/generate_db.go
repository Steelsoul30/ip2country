package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"ip2country/internal/dbgenerator"
)

var createCmd = &cobra.Command{
	Use:     "create-db",
	Aliases: []string{"gdb", "cdb", "create-database", "generate-database"},
	Short:   "Create database",
	Long:    "Generate database to serialize the database schema to a file in preparation for runtime use.",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("zippath")
		slog.Info(fmt.Sprintf("Creating database...\n flag is %s\n", path))
		absPath, err := filepath.Abs(path)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to get absolute path from path %s: %v\n", path, err))
			return
		}
		if stat, err := os.Stat(absPath); os.IsNotExist(err) || stat == nil {
			slog.Error(fmt.Sprintf("The file %s does not exist: %v\n", absPath, err))
			return
		}

		slog.Info(fmt.Sprintf("Creating database from zip file %s\n", absPath))
		dbGen := dbgenerator.NewDbGenerator()
		_, err = dbGen.DirectFromZip(absPath)
		if err != nil {
			slog.Error(fmt.Sprintf("Error creating database: %v\n", err))
			return
		}
		err = dbGen.SaveInfo(filepath.Dir(absPath) + "/geodata.dat")
		if err != nil {
			slog.Error(fmt.Sprintf("Error saving database: %v\n", err))
			return
		}
		slog.Info("Database created successfully")
	},
}

func init() {
	createCmd.Flags().StringP("zippath", "p", "db/geolite2.zip", "Path to the zip file containing the database")
	rootCmd.AddCommand(createCmd)
}
