package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/cli/go-gh/v2"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/browser"
)

func main() {
	currentBranch := getCurrentBranch()
	repositoryUrl := getRepositoryUrl()
	prUrl := getPrUrl(repositoryUrl, currentBranch)

	browser.OpenURL(prUrl)
}

func getCurrentBranch() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		log.Fatal(err)
	}

	head, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}

	ref := head.Name().String()
	return strings.Replace(ref, "refs/heads/", "", 1)
}

func getRepositoryUrl() string {
	stdOut, _, err := gh.Exec("browse", "-n")
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSuffix(stdOut.String(), "\n")
}

func getPrUrl(repository string, branch string) string {
	prUrl, err := url.JoinPath(repository, "pull", branch)
	if err != nil {
		log.Fatal(err)
	}

	return prUrl
}
