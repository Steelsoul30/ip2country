package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the service",
	Long:  "Run the service. Starts a server that listens to incoming requests",
	Run: func(cmd *cobra.Command, args []string) {
		s := "gopher"
		fmt.Printf("Hello and welcome, %s!\n", s)
		for i := 1; i <= 5; i++ {
			fmt.Println("i =", 100/i)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
