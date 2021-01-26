package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	user    string
	rootCmd = &cobra.Command{
		Use:   "github-cli",
		Short: "github-cli is a command line interface for use with multiple github repositories",
	}
)

func init() {
	rootCmd.AddCommand(whichCmd)
	rootCmd.AddCommand(mainCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
