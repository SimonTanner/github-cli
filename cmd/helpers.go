package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
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
