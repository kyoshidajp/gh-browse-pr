package main

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

func main() {
	currentBranch := getCurrentBranch()
	repositoryUrl := getRepositoryUrl()
	prUrl := GetPrUrl(repositoryUrl, currentBranch)

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

func GetNewPrUrl(repository string, branch string) string {
	prUrl, err := url.JoinPath(repository, "compare", branch)
	if err != nil {
		log.Fatal(err)
	}

	parsedUrl, err := url.Parse(prUrl)
	if err != nil {
		log.Fatal(err)
	}

	queries := parsedUrl.Query()
	queries.Add("expand", "1")
	parsedUrl.RawQuery = queries.Encode()
	return parsedUrl.String()
}

func GetPrUrl(repository string, branch string) string {
	branchOrNumber := ""

	if IsNumberString(branch) {
		repo := strings.Replace(repository, "https://github.com/", "", 1)
		repoQuery := fmt.Sprintf("--repo=%s", repo)
		headQuery := fmt.Sprintf("--head=%s", branch)
		searchedPrNumber, _, err := gh.Exec("search", "prs", repoQuery, headQuery, "--json=number", "-q=.[].number")
		if err != nil {
			return GetNewPrUrl(repository, branch)
		}

		prNumber := strings.TrimSuffix(searchedPrNumber.String(), "\n")
		if prNumber == "" {
			return GetNewPrUrl(repository, branch)
		} else {
			branchOrNumber = prNumber
		}
	} else {
		branchOrNumber = branch
	}

	prUrl, err := url.JoinPath(repository, "pull", branchOrNumber)
	if err != nil {
		log.Fatal(err)
	}

	return prUrl
}

func IsNumberString(value string) bool {
	if _, err := strconv.Atoi(value); err == nil {
		return true
	}

	return false
}
