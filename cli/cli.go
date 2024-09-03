package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "binp",
	Short: "A cli tool for the binp pastebin service",
	Long:  "A cli tool for the binp pastebin service",
}

func Execute() error {
	return rootCmd.Execute()
}
