package git

import (
	"context"
	"fmt"
	"log"
)

type GitPullRequest struct {
	Number     uint
	Title      string
	URL        string
	State      string
	LastCommit string
}

type GitPullRequestFileChanges struct {
	Filename    string `json:"filename"`
	Patch       string `json:"patch"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	Changes     int    `json:"changes"`
	FileContent string `json:"file_content"`
}

type GitClient interface {
	ListPullRequests(ctx context.Context, repoName, repoOwner string, userID uint) ([]*GitPullRequest, error)
	FetchFileDiffs(ctx context.Context, repoName, repoOwner string, prNumber uint, userID uint) ([]*GitPullRequestFileChanges, error)
}

func NewGitClient(ctx context.Context, clientType string) (GitClient, error) {
	switch clientType {
	case "github":
		client, err := NewGithubClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create GitHub client: %v", err)
			return nil, err
		}
		return client, nil
	case "azure":
		client, err := NewAzureDevopsClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create Azure Devops client: %v", err)
			return nil, err
		}
		return client, nil
	default:
		log.Fatalf("Unknown client type: %s", clientType)
		return nil, fmt.Errorf("unknown client type: %s", clientType)
	}
}
