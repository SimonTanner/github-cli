package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	us "github.com/SimonTanner/github-cli/user"
)

// getValueFromConfigList returns the value for a given key
func getValueFromConfigList(list bytes.Buffer, key string) (string, error) {
	listStr := list.String()
	var val string
	err := fmt.Errorf("unable to get %s from config", key)

	idx := strings.Index(listStr, key)
	if idx == -1 {
		return val, err
	}
	valLine := strings.Split(listStr[idx:], "\n")[0]
	if len(valLine) == 0 {
		return val, err
	}
	val = strings.Split(valLine, "=")[1]

	return val, nil
}

func setConfigVar(varName, value string) error {
	cmd := exec.Command("git", "config", "--local", varName, fmt.Sprintf("\"%s\"", value))

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func setUserProfile(user string, add, force bool) error {
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

	err = setUser(currentUser, true)
	if err != nil {
		return err
	}

	fmt.Println("Git user successfully set to", user)
	return nil
}

func setUser(user us.User, output bool) error {
	nameErr := setConfigVar(userNameKey, user.Name)
	if nameErr != nil {
		return nameErr
	}
	if output {
		fmt.Printf("%s set to %s, ", userNameKey, user.Name)
	}

	emailErr := setConfigVar(userEmailKey, user.Email)
	if emailErr != nil {
		return emailErr
	}
	if output {
		fmt.Printf("%s set to %s\n", userEmailKey, user.Email)
	}

	return nil
}

func inputNewUser(user string, overwrite bool) (us.User, error) {
	var err error
	reader := bufio.NewReader(os.Stdin)
	delim := byte('\n')

	fmt.Print("Please enter the user name: ")
	uName, err := reader.ReadString(delim)
	if err != nil {
		return us.User{}, err
	}
	uName = strings.ReplaceAll(uName, string(delim), "")

	fmt.Print("Please enter the email address: ")
	uEmail, err := reader.ReadString(delim)
	if err != nil {
		return us.User{}, err
	}
	uEmail = strings.ReplaceAll(uEmail, string(delim), "")

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

func inputAccessToken(profile string, user us.User) (us.User, error) {
	var err error
	reader := bufio.NewReader(os.Stdin)
	delim := byte('\n')

	fmt.Print("please enter your github access token: ")
	token, err := reader.ReadString(delim)
	if err != nil {
		return us.User{}, err
	}
	user.AccessToken = strings.ReplaceAll(token, string(delim), "")
	_, u, err := us.SaveUser(profile, user, true)
	if err != nil {
		return us.User{}, err
	}

	return u, nil
}
