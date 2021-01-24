package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
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
	userName := "user.name"

	cmd := exec.Command("git", "config", "--list", "--local")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	user, err := getValueFromConfigList(out, userName)
	if err != nil || user == "" {
		fmt.Println("No user set in local config, checking global user")
		cmd := exec.Command("git", "config", "--list", "--global")
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	user, err = getValueFromConfigList(out, userName)
	if err != nil || user == "" {
		return unknownUserErr
	}

	fmt.Println("Active user:", user)
	return nil
}

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
