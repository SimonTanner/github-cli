package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	workUser = "work"
)

var (
	workCmd = &cobra.Command{
		Use:   "work",
		Short: fmt.Sprintf("set git user to \"%s\" locally", workUser),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := work(workUser, add, force); err != nil {
				return err
			}
			return nil
		},
		SilenceUsage: true,
	}
)

func init() {
	workCmd.Flags().BoolVarP(&add, "add", "a", false, "use to add new user profile details")
	workCmd.Flags().BoolVarP(&force, "force", "f", false, "use to overwrite user profile details")
}

func work(user string, add, force bool) error {
	err := setUserProfile(user, add, force)
	return err
}
