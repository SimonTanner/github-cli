package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	us "github.com/SimonTanner/github-cli/user"

	"github.com/spf13/cobra"
)

const (
	userEmailKey = "user.email"
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
	mainCmd.MarkFlagRequired("user")
}

func main(user string) error {
	branchCmd := exec.Command("git", "branch")
	var (
		outBranch, errBranch bytes.Buffer
	)

	branchCmd.Stdout = &outBranch
	branchCmd.Stderr = &errBranch
	branchErr := branchCmd.Run()
	if branchErr != nil {
		fmt.Println("Not a git repository. If you wish to initialise a repository try \"github-cli set [private]\"")
		return nil
	}

	fmt.Println("Current branch:", strings.Trim(outBranch.String(), "/n* "))

	u, uErr := us.GetUser(user)
	if uErr != nil {
		return fmt.Errorf("error getting user: %w", uErr)
	}

	nameErr := setConfigVar(userNameKey, u.Name)
	fmt.Println(fmt.Sprintf("setting %s to %s", userNameKey, u.Name))
	if nameErr != nil {
		return nameErr
	}
	emailErr := setConfigVar(userEmailKey, u.Email)
	if emailErr != nil {
		return emailErr
	}
	fmt.Println("User set to", user)
	return nil
}
