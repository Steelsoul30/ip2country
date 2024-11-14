package cmd

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	"ip2country/internal/config"
	"ip2country/internal/logger"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "ip2country",
	Short: "ip2country is a service that provides country information based on IP addresses",
	Long:  `ip2country is a service that provides country information based on IP addresses.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Running command: %s\n", cmd.Name())
		localConfig, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}
		cfg = localConfig
		logger.InitLogger(cfg)
		slog.Info("Logger initialized")

		config.PrintConfigToLog(cfg, "")
		return nil
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing ip2country: %v", err)
	}
}
