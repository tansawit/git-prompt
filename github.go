package main

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v29/github" // with go modules enabled (GO111MODULE=on or outside GOPATH)
)

func getGitHubClient(ghToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}

// GetUserRepos returns a slice of the user's public GitHub Repositories
// Implementation Credit: https://github.com/lox/alfred-github-jump/repos.go
func githubGetUserRepos(ghToken string) ([]*github.Repository, map[string]github.Repository, error) {
	fmt.Println("Fetching latest list of GitHub repos")
	client := getGitHubClient(ghToken)
	ctx := context.Background()
	var repoMap = map[string]github.Repository{}

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 45},
		Sort:        "pushed",
	}

	repos := []*github.Repository{}

	for {
		result, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			return repos, repoMap, err
		}
		repos = append(repos, result...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	for _, repo := range repos {
		repoMap[*repo.Name] = *repo
	}

	return repos, repoMap, nil
}
