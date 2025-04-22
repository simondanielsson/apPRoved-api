package git

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	adogit "github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
)

const placeholderBinaryContent string = "Binary file - not displayed"

type AzureDevopsClient struct {
	client adogit.Client
	mutex  *sync.Mutex
}

func NewAzureDevopsClient(ctx context.Context) (*AzureDevopsClient, error) {
	organizationUrl := os.Getenv("AZURE_DEVOPS_ORG_URL")
	if organizationUrl == "" {
		return nil, fmt.Errorf("AZURE_DEVOPS_ORG_URL environment variable is not set")
	}
	personalAccessToken := os.Getenv("AZURE_DEVOPS_PAT")
	if personalAccessToken == "" {
		return nil, fmt.Errorf("AZURE_DEVOPS_PAT environment variable is not set")
	}

	connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
	client, err := adogit.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure DevOps client: %w", err)
	}
	log.Printf("Initialized Azure DevOps client")
	azureClient := &AzureDevopsClient{
		client: client,
		mutex:  &sync.Mutex{},
	}
	return azureClient, nil
}

// ListPullRequests lists all pull requests for a given repository [repoName] in project [repoOwner].
func (c *AzureDevopsClient) ListPullRequests(ctx context.Context, repoName, repoOwner string, userID uint) ([]*GitPullRequest, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	project := repoOwner
	searchArgs := adogit.GetPullRequestsArgs{
		Project:        &project,
		RepositoryId:   &repoName,
		SearchCriteria: &adogit.GitPullRequestSearchCriteria{},
	}
	response, err := c.client.GetPullRequests(ctx, searchArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to list pull requests: %w", err)
	}

	pullRequests := make([]*GitPullRequest, 0, len(*response))
	log.Printf("Found %d pull requests for repo %s", len(*response), repoName)
	for _, pr := range *response {
		var pullRequest GitPullRequest
		if pr.PullRequestId == nil || *pr.PullRequestId < 0 {
			log.Printf("Invalid pull request ID: %d, skipping", *pr.PullRequestId)
			continue
		} else {
			pullRequest.Number = uint(*pr.PullRequestId)
		}
		if pr.Title != nil {
			pullRequest.Title = *pr.Title
		}
		if pr.Url != nil {
			pullRequest.URL = *pr.Url
		}
		if pr.Status != nil {
			pullRequest.State = string(*pr.Status)
		}
		if pr.LastMergeCommit != nil && pr.LastMergeCommit.CommitId != nil {
			pullRequest.LastCommit = *pr.LastMergeCommit.CommitId
		}
		pullRequests = append(pullRequests, &pullRequest)
	}
	return pullRequests, nil
}

func (c *AzureDevopsClient) FetchFileDiffs(ctx context.Context, repoName, repoOwner string, prNumber uint, userID uint) ([]*GitPullRequestFileChanges, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// We need to see the target branch of the PR
	prNumberInt := int(prNumber)
	pr, err := c.client.GetPullRequest(ctx, adogit.GetPullRequestArgs{
		Project:       &repoOwner,
		RepositoryId:  &repoName,
		PullRequestId: &prNumberInt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull request details: %w", err)
	}

	// Find the (last) iteration (commit) for which to find the patch for
	iterations, err := c.client.GetPullRequestIterations(ctx, adogit.GetPullRequestIterationsArgs{
		Project:       &repoOwner,
		RepositoryId:  &repoName,
		PullRequestId: &prNumberInt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull request iterations: %w", err)
	}
	var lastIterationId int
	if iterations == nil || len(*iterations) == 0 {
		log.Printf("No iterations found for PR %d", prNumber)
		lastIterationId = 1
	} else {
		lastIterationId = len(*iterations)
	}

	compareTo := 0
	iterationChanges, err := c.client.GetPullRequestIterationChanges(ctx, adogit.GetPullRequestIterationChangesArgs{
		Project:       &repoOwner,
		RepositoryId:  &repoName,
		PullRequestId: &prNumberInt,
		IterationId:   &lastIterationId,
		CompareTo:     &compareTo,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file diffs: %w", err)
	}
	if iterationChanges.ChangeEntries == nil {
		log.Printf("No file changes found for PR %d", prNumber)
		return []*GitPullRequestFileChanges{}, nil
	}

	fileChanges := make([]*GitPullRequestFileChanges, 0, len(*iterationChanges.ChangeEntries))
	for _, change := range *iterationChanges.ChangeEntries {
		if change.Item == nil {
			log.Printf("Item is nil for change entry, skipping")
			continue
		}

		// The Item is actually a map[string]interface{}, not a GitItem struct
		itemMap, ok := change.Item.(map[string]interface{})
		if !ok {
			log.Printf("Failed to assert change.Item as map[string]interface{}, skipping")
			continue
		}

		// Extract the path from the map
		pathInterface, ok := itemMap["path"]
		if !ok {
			log.Printf("Path key doesn't exist in item map, skipping")
			continue
		}

		path, ok := pathInterface.(string)
		if !ok {
			log.Printf("Path is not a string, skipping")
			continue
		}

		targetBranch := cleanBranchName(*pr.TargetRefName)
		sourceBranch := cleanBranchName(*pr.SourceRefName)
		getNewContentArgs := adogit.GetItemArgs{
			RepositoryId: &repoName,
			Project:      &repoOwner,
			Path:         &path,
			VersionDescriptor: &adogit.GitVersionDescriptor{
				Version:     &sourceBranch,
				VersionType: &adogit.GitVersionTypeValues.Branch,
			},
			IncludeContentMetadata: boolPtr(true),
			IncludeContent:         boolPtr(true),
		}
		newContent, newContentType, err := c.readFileContent(ctx, prNumberInt, getNewContentArgs)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// Azure does not track number of lines added, changed, or deleted so we default to the zero-value for these fields
		fileChange := GitPullRequestFileChanges{
			Filename:    path,
			Additions:   0,
			Deletions:   0,
			Changes:     0,
			FileContent: newContent,
		}

		switch *change.ChangeType {
		case adogit.VersionControlChangeTypeValues.Add:
			if newContentType == "base64Encoded" {
				fileChange.Patch = createPatch(path, path, "", placeholderBinaryContent)
			} else {
				fileChange.Patch = createPatch(path, path, "", newContent)
			}

		case adogit.VersionControlChangeTypeValues.Edit, adogit.VersionControlChangeTypeValues.Delete:
			if newContentType == "base64Encoded" {
				fileChange.Patch = placeholderBinaryContent
			} else {
				getOldItemArgs := adogit.GetItemArgs{
					RepositoryId: &repoName,
					Project:      &repoOwner,
					Path:         &path,
					VersionDescriptor: &adogit.GitVersionDescriptor{
						Version:     &targetBranch,
						VersionType: &adogit.GitVersionTypeValues.Branch,
					},
					IncludeContentMetadata: boolPtr(true),
					IncludeContent:         boolPtr(true),
				}
				oldContent, oldContentType, err := c.readFileContent(ctx, prNumberInt, getOldItemArgs)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				if oldContentType == "base64Encoded" {
					oldContent = placeholderBinaryContent
				}
				fileChange.Patch = createPatch(path, path, oldContent, newContent)
			}
		default:
			log.Printf("Unknown change type for %s: %s", path, *change.ChangeType)
			continue
		}

		fileChanges = append(fileChanges, &fileChange)
	}
	return fileChanges, nil
}

// readFileContent reads the content of a file in a pull request.
func (c *AzureDevopsClient) readFileContent(ctx context.Context, prNumber int, args adogit.GetItemArgs) (string, string, error) {
	item, err := c.client.GetItem(ctx, args)
	if err != nil {
		return "", "", fmt.Errorf("failed to get item %s for PR %d: %w", *args.Path, prNumber, err)
	}
	if item.Content == nil {
		return "", "", fmt.Errorf("no content found for PR %d", prNumber)
	}
	if item.ContentMetadata.ContentType == nil {
		return "", "", fmt.Errorf("no content type found for PR %d", prNumber)
	}
	return *item.Content, *item.ContentMetadata.ContentType, nil
}

// createPatch creates a patch string from the old and new content.
func createPatch(oldName, newName, oldContent, newContent string) string {
	edits := myers.ComputeEdits(span.URIFromPath(oldName), oldContent, newContent)
	return fmt.Sprint(gotextdiff.ToUnified(oldName, newName, oldContent, edits))
}

// cleanBranchName cleans a branch name.
func cleanBranchName(branchName string) string {
	prefix := "refs/heads/"
	if len(branchName) > len(prefix) && branchName[:len(prefix)] == prefix {
		return branchName[len(prefix):]
	}
	return branchName
}

// boolPtr converts a boolean value to a boolean pointer.
func boolPtr(b bool) *bool {
	return &b
}
