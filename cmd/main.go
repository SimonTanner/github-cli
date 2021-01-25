package cmd

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	mainCmd = &cobra.Command{
		Use:   "main",
		Short: "set the main git user locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := main(user); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	mainCmd.Flags().StringVarP(&user, "user", "u", "", "the user you wish to set locally")
}

func main(user string) error {
	cmd := exec.Command("git", "config", "--local", "user.name", user)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
