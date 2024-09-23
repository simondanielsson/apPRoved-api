package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/go-github/v64/github"
	"golang.org/x/oauth2"
)

type GithubPullRequest struct {
	Number     uint
	Title      string
	URL        string
	State      string
	LastCommit string
}

type GithubPullRequestFileChanges struct {
	Filename  string `json:"filename"`
	Patch     string `json:"patch"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
}

type GithubClient struct {
	client *github.Client
	mutex  *sync.Mutex
}

func NewGithubClient(ctx context.Context) (*GithubClient, error) {
	client := GithubClient{
		client: nil,
		mutex:  &sync.Mutex{},
	}

	if err := client.connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize GitHub client")
	}

	return &client, nil
}

func (c *GithubClient) connect(ctx context.Context) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatalf("GITHUB_TOKEN is not set")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	if tc == nil {
		log.Fatal("Failed to initialize OAuth2 client")
	}

	client := github.NewClient(tc)
	if client == nil {
		log.Fatalf("Failed to initialize GitHub client")
	}
	c.client = client

	log.Printf("Initialized GitHub client")
	return nil
}

func (c *GithubClient) ListPullRequests(ctx context.Context, repoName, repoOwner string, userID uint) ([]*GithubPullRequest, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	log.Printf("Fetching PRs for %s/%s", repoOwner, repoName)
	// fetch PR's by page
	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{Page: 0, PerPage: 30},
	}

	var prs []*GithubPullRequest
	for {
		fetchedPRs, resp, err := c.client.PullRequests.List(ctx, repoOwner, repoName, opts)
		if err != nil {
			if r, _ := err.(*github.ErrorResponse); r.Response.StatusCode == 404 {
				fetchedPRs = []*github.PullRequest{}
			} else {
				return nil, err
			}
		}

		for _, pr := range fetchedPRs {
			number := *pr.Number
			if number < 0 {
				number = 0
			}
			uint_number := uint(number)

			commits, _, err := c.client.PullRequests.ListCommits(ctx, repoOwner, repoName, *pr.Number, nil)
			if err != nil {
				return nil, err
			}

			prs = append(prs, &GithubPullRequest{
				Number:     uint_number,
				Title:      *pr.Title,
				URL:        *pr.URL,
				State:      *pr.State,
				LastCommit: *commits[len(commits)-1].SHA,
			})
			log.Printf("#%d %s (%s)\n", *pr.Number, *pr.Title, *pr.URL)
		}

		// If there are no more pages, break the loop (0 is the zero value for NextPage))
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	log.Printf("Fetched %d PRs for %s/%s", len(prs), repoOwner, repoName)
	return prs, nil
}

func (c *GithubClient) FetchFileDiffs(ctx context.Context, repoName, repoOwner string, prNumber uint, userID uint) ([]*GithubPullRequestFileChanges, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	files, _, err := c.client.PullRequests.ListFiles(ctx, repoOwner, repoName, int(prNumber), nil)
	if err != nil {
		return nil, err
	}

	var fc []*GithubPullRequestFileChanges
	for _, file := range files {
		diff := &GithubPullRequestFileChanges{
			Filename:  *file.Filename,
			Patch:     SafeString(file.Patch, "Cannot display patch for binary file"),
			Additions: *file.Additions,
			Deletions: *file.Deletions,
			Changes:   *file.Changes,
		}
		fc = append(fc, diff)
	}

	return fc, nil
}
