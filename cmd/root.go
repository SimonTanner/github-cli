package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	add     bool
	force   bool
	private bool
	rootCmd = &cobra.Command{
		Use:   "github-cli",
		Short: "github-cli is a command line interface for use with multiple github repositories",
	}
)

func init() {
	rootCmd.AddCommand(whichCmd)
	rootCmd.AddCommand(mainCmd)
	rootCmd.AddCommand(workCmd)
	rootCmd.AddCommand(setCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
