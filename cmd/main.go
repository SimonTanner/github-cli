package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	us "github.com/SimonTanner/github-cli/user"

	"github.com/spf13/cobra"
)

const (
	userEmailKey = "user.email"
	mainUser     = "main"
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
		SilenceUsage: true,
	}
)

func init() {
	mainCmd.Flags().StringVarP(&user, "user", "u", "", "the user you wish to set locally")
	// mainCmd.MarkFlagRequired("user")
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

	if user == "" {
		user = mainUser
		fmt.Println("user is main")
	}

	fmt.Println("Current branch:", strings.Trim(outBranch.String(), "\n* "))

	u, uErr := us.GetUser(user)
	// if uErr != nil && uErr == fmt.Errorf("no user found with profile name %s", user) {
	if uErr != nil {
		var err error
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Please enter the user name: ")
		uName, _ := reader.ReadString('\n')
		fmt.Print("Please enter the email address: ")
		uEmail, _ := reader.ReadString('\n')
		newUser := us.User{
			Name:  uName,
			Email: uEmail,
		}

		_, u, err = us.SaveUser(user, newUser)
		if err != nil {
			fmt.Println(err)
		}
	}

	nameErr := setConfigVar(userNameKey, u.Name)
	fmt.Printf("setting %s to %s", userNameKey, u.Name)
	if nameErr != nil {
		return nameErr
	}
	emailErr := setConfigVar(userEmailKey, u.Email)
	fmt.Printf("setting %s to %s", userEmailKey, u.Email)
	if emailErr != nil {
		return emailErr
	}
	fmt.Println("User successfully set to", user)
	return nil
}
