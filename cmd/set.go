package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	us "github.com/SimonTanner/github-cli/user"

	"github.com/google/go-github/v33/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
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

	ErrRepoCreationFailed = errors.New("failed to create remote repository")
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

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println("Current Directory:", currDir)
	dirSlice := strings.Split(currDir, "/")
	repoName := dirSlice[len(dirSlice)-1]

	newRepo, statusCode, err := createRemoteRepo(repoName, private)
	if err != nil {
		return err
	}

	if newRepo == nil {
		return ErrRepoCreationFailed
	}
	fmt.Printf("%s repository successfully created, http response status code: %d", repoName, statusCode)
	fmt.Printf("\n%+v", *newRepo)

	return nil
}

func createRemoteRepo(name string, private bool) (*github.Repository, int, error) {
	ctx := context.Background()
	st := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: "5f4a6a26baf2c8a2c8af1a69970a24ab0af77de5",
		},
	)

	oaClient := oauth2.NewClient(ctx, st)

	gitClient := github.NewClient(oaClient)

	repo := github.Repository{
		Name:    &name,
		Private: &private,
	}
	org := ""

	newRepo, resp, err := gitClient.Repositories.Create(ctx, org, &repo)
	if err != nil {
		return newRepo, resp.StatusCode, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return newRepo, resp.StatusCode, fmt.Errorf("error creating repository, server returned status code: %d", resp.StatusCode)
	}

	return newRepo, resp.StatusCode, nil
}
