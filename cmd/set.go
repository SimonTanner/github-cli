package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	us "github.com/SimonTanner/github-cli/user"
	"github.com/spf13/cobra"
)

var (
	setCmd = &cobra.Command{
		Use:   "set",
		Short: "initialises a github repository locally and makes it private if [private] flag set",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := set(private); err != nil {
				return err
			}
			return nil
		},
		SilenceUsage: true,
	}
)

func init() {
	setCmd.Flags().BoolVarP(&private, "private", "p", false, "use to make repository private")
}

func set(private bool) error {
	cmd := exec.Command("git", "init")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmdErr := cmd.Run()
	if cmdErr != nil {
		fmt.Println("Bugger")
		log.Fatal(cmdErr)
		return cmdErr
	}

	user, err := us.GetUser(mainUser)
	if err != nil {
		return err
	}
	err = setUser(user, false)
	if err != nil {
		return err
	}

	return nil
}
