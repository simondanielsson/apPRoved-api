package utils

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v64/github"
)

var ghClient map[uint]*github.Client

func initGithubClient(userID uint) *github.Client {
	gh_token := os.Getenv("GITHUB_TOKEN")
	ghClient[userID] = github.NewClient(nil).WithAuthToken(gh_token)
	log.Printf("Initialized GitHub client for user %d", userID)
	return ghClient[userID]
}

func getOrCreateGithubClient(userID uint) *github.Client {
	client := ghClient[userID]
	if client == nil {
		log.Printf("GitHub client not initialized for user %d. Initializing.", userID)
		client := initGithubClient(userID)
		return client
	}
	return client
}

type GithubPullRequest struct {
	Number uint
	Title  string
	URL    string
}

func ListPullRequests(ctx context.Context, repoName, repoOwner string, userID uint) ([]*GithubPullRequest, error) {
	client := getOrCreateGithubClient(userID)
	fetched_prs, _, err := client.PullRequests.List(ctx, repoOwner, repoName, &github.PullRequestListOptions{
		State: "open",
	})
	if err != nil {
		return nil, err
	}
	fetched_prs[0].GetURL()

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
