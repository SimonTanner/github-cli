package cmd

import (
	"bufio"
	"bytes"
	"errors"
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
	add     bool
	force   bool
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
	branchCmd := exec.Command("git", "branch")
	var (
		outBranch, errBranch bytes.Buffer
		currentUser          us.User
		err                  error
	)

	branchCmd.Stdout = &outBranch
	branchCmd.Stderr = &errBranch
	branchErr := branchCmd.Run()
	if branchErr != nil {
		fmt.Println("Not a git repository. If you wish to initialise a repository try \"github-cli set [private]\"")
		return nil
	}

	fmt.Println("Current branch:", strings.Trim(outBranch.String(), "\n* "))

	if add && !force {
		currentUser, err = inputNewUser(user, false)
		if err != nil {
			if errors.Is(err, us.ErrUserExists) {
				fmt.Printf("user profile \"%s\" already exists\n", user)
				fmt.Println("to overwrite this enter github-cli -a -f")
				return nil
			}
			return err
		}
	} else if add && force {
		currentUser, err = inputNewUser(user, true)
		if err != nil {
			return err
		}
	} else {
		currentUser, err = us.GetUser(user)
		if err != nil {
			if errors.Is(err, us.ErrNoUserFound) {
				fmt.Printf("user profile \"%s\" not currently set\n", user)
				fmt.Println("To add a new user use enter -a")
				return nil
			}
			return err
		}
	}

	nameErr := setConfigVar(userNameKey, currentUser.Name)
	fmt.Printf("Setting %s to %s", userNameKey, currentUser.Name)
	if nameErr != nil {
		return nameErr
	}
	emailErr := setConfigVar(userEmailKey, currentUser.Email)
	fmt.Printf("Setting %s to %s", userEmailKey, currentUser.Email)
	if emailErr != nil {
		return emailErr
	}

	fmt.Println("Git user successfully set to", mainUser)
	return nil
}

func inputNewUser(user string, overwrite bool) (us.User, error) {
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

	_, u, err := us.SaveUser(user, newUser, overwrite)
	if err != nil {
		return us.User{}, err
	}
	return u, nil
}
