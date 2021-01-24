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

	cmd := exec.Command("git", "config", "--list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	outStr := out.String()

	idx := strings.Index(outStr, userName)
	if idx == 0 {
		return unknownUserErr
	}
	userLine := strings.Split(outStr[idx:], "\n")[0]
	if len(userLine) == 0 {
		return unknownUserErr
	}
	user := strings.Split(userLine, "=")[1]
	fmt.Println("Active user:", user)

	return nil
}
