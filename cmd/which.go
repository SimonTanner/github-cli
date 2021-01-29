package cmd

import (
	"bytes"
	"errors"
	"fmt"
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
	var (
		out                 bytes.Buffer
		userName, userEmail string
	)

	cmd := exec.Command("git", "config", "--list", "--local")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		out, err = checkGlobalConf()
		if err != nil {
			fmt.Println("No user found in global config")
			return err
		}
	}
	userName, err = getValueFromConfigList(out, userNameKey)
	if err != nil {
		return err
	}

	userEmail, err = getValueFromConfigList(out, userEmailKey)
	if err != nil {
		return err
	}

	fmt.Printf("Active user: %s, email: %s", userName, userEmail)
	return nil
}

func checkGlobalConf() (bytes.Buffer, error) {
	fmt.Println("No user set in local config, checking global user")
	cmd := exec.Command("git", "config", "--list", "--global")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return out, err
	}

	return out, nil
}
