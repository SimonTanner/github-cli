package cmd

import (
	"github.com/spf13/cobra"
)

const (
	userEmailKey = "user.email"
	mainUser     = "main"
)

var (
	mainCmd = &cobra.Command{
		Use:   "main",
		Short: "set git user to \"main\" locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := main(mainUser, add, force); err != nil {
				return err
			}
			return nil
		},
		SilenceUsage: true,
	}
)

func init() {
	mainCmd.Flags().BoolVarP(&add, "add", "a", false, "use to add new user profile details")
	mainCmd.Flags().BoolVarP(&force, "force", "f", false, "use to overwrite user profile details")
}

func main(user string, add, force bool) error {
	err := setUserProfile(user, add, force)
	return err
}
