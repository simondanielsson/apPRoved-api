package utils

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v64/github"
)

func initGithubClient(userID uint) *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatalf("GITHUB_TOKEN is not set")
	}
	client := github.NewClient(nil).WithAuthToken(token)
	log.Printf("Initialized GitHub client for user %d", userID)
	return client
}

type GithubPullRequest struct {
	Number uint
	Title  string
	URL    string
}

func ListPullRequests(ctx context.Context, repoName, repoOwner string, userID uint) ([]*GithubPullRequest, error) {
	client := initGithubClient(userID)
	fetched_prs, _, err := client.PullRequests.List(ctx, repoOwner, repoName, &github.PullRequestListOptions{
		State: "open",
	})
	log.Printf("Fetched %d PRs, checking error: %s", len(fetched_prs), err.Error())
	if err != nil {
		return nil, err
	}

	prs := make([]*GithubPullRequest, len(fetched_prs))
	for _, pr := range fetched_prs {
		number := pr.GetNumber()
		if number < 0 {
			number = 0
		}
		uint_number := uint(number)

		prs = append(prs, &GithubPullRequest{
			Number: uint_number,
			Title:  pr.GetURL(),
			URL:    pr.GetTitle(),
		})
	}

	return prs, nil
}
