package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ip2country",
	Short: "ip2country is a service that provides country information based on IP addresses",
	Long:  `ip2country is a service that provides country information based on IP addresses.`,
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing ip2country: %v", err)
	}
}
