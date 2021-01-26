package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	userNameKey = "user.name"
)

var (
	whichCmd = &cobra.Command{
		Use:   "which",
		Short: "which config variables",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := which(); err != nil {
				return err
			}
			return nil
		},
	}

	unknownUserErr = errors.New("no user name found in git config")
)

func which() error {
	cmd := exec.Command("git", "config", "--list", "--local")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	user, err := getValueFromConfigList(out, userNameKey)
	flag := "--local"
	if err != nil || user == "" {
		flag = "--global"
		fmt.Println("No user set in local config, checking global user")
		cmd := exec.Command("git", "config", "--list", flag)
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	user, err = getValueFromConfigList(out, userNameKey)
	if err != nil || user == "" {
		return unknownUserErr
	}

	fmt.Println("Active user:", user)
	return nil
}
