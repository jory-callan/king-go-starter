package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("MyApp Version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
