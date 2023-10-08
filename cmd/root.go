package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/cli/go-gh/v2"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/browser"
)

const (
	GitHubURL = "https://github.com"
)

func Execute() {
	branch := currentBranch()
	repositoryURL := repositoryURL()
	prURL := PrURL(repositoryURL, branch)

	if !strings.HasPrefix(prURL, GitHubURL) {
		message := fmt.Sprintf("Can't open URL: %s", prURL)
		log.Fatal(message)
	}

	browser.OpenURL(prURL)
}

func currentBranch() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
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

func repositoryURL() string {
	stdOut, _, err := gh.Exec("browse", "-n")
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSuffix(stdOut.String(), "\n")
}

func NewPrURL(repository string, branch string) string {
	prURL, err := url.JoinPath(repository, "compare", branch)
	if err != nil {
		log.Fatal(err)
	}

	parsedURL, err := url.Parse(prURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := parsedURL.Query()
	queries.Add("expand", "1")
	parsedURL.RawQuery = queries.Encode()
	return parsedURL.String()
}

func PrURL(repository string, branch string) string {
	branchOrNumber := ""

	if IsNumberString(branch) {
		repo := strings.Replace(repository, GitHubURL+"/", "", 1)
		repoQuery := fmt.Sprintf("--repo=%s", repo)
		headQuery := fmt.Sprintf("--head=%s", branch)
		searchedPrNumber, _, err := gh.Exec("search", "prs", repoQuery, headQuery, "--json=number", "-q=.[].number")
		if err != nil {
			return NewPrURL(repository, branch)
		}

		prNumber := strings.TrimSuffix(searchedPrNumber.String(), "\n")
		if prNumber == "" {
			return NewPrURL(repository, branch)
		} else {
			branchOrNumber = prNumber
		}
	} else {
		branchOrNumber = branch
	}

	prURL, err := url.JoinPath(repository, "pull", branchOrNumber)
	if err != nil {
		log.Fatal(err)
	}

	return prURL
}

func IsNumberString(value string) bool {
	if _, err := strconv.Atoi(value); err == nil {
		return true
	}

	return false
}
