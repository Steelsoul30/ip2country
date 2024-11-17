package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"ip2country/internal/config"
	"ip2country/internal/ip2country/handler"
	"ip2country/internal/ip2country/store"
	"ip2country/internal/logger"
	"ip2country/internal/router"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the service",
	Long:  "Run the service. Starts a server that listens to incoming requests",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Running command: %s\n", cmd.Name())
		localConfig, err := config.LoadConfig()
		if err != nil {
			return
		}
		cfg = localConfig
		logger.InitLogger(cfg)
		slog.Info("Logger initialized")
		config.PrintConfigToLog(cfg, "")
		slog.Info("Initializing data store")
		storeImpl, err := store.NewStore(cfg, cmd)
		if err != nil {
			return
		}
		slog.Info("Data store initialized")
		handler.SetStore(storeImpl)

		slog.Info(fmt.Sprintf("Starting server on %d", cfg.Port))
		router.StartServer(cfg)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
