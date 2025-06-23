package github

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v72/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	orgs   []string
}

func NewClient(orgs []string) (*Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is required")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &Client{
		client: client,
		orgs:   orgs,
	}, nil
}

func (g *Client) GetIssues(
	ctx context.Context,
	since, until time.Time,
	author string,
	includeBody bool,
) ([]Issue, error) {
	var allIssues []Issue

	// Search for issues created by the author
	query := fmt.Sprintf("author:%s created:%s..%s type:issue",
		author,
		since.Format("2006-01-02"),
		until.Format("2006-01-02"))

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		result, resp, err := g.client.Search.Issues(ctx, query, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to search issues: %w", err)
		}

		for _, issue := range result.Issues {
			// Extract repository info from the issue URL
			if issue.HTMLURL == nil {
				continue
			}

			// Parse repository from URL like https://github.com/owner/repo/issues/123
			urlParts := strings.Split(*issue.HTMLURL, "/")
			if len(urlParts) < 6 {
				continue
			}

			repoFullName := urlParts[3] + "/" + urlParts[4]

			// Check if the repository belongs to target organizations (if specified)
			shouldInclude := len(g.orgs) == 0 // Include all if no orgs specified
			if !shouldInclude {
				for _, org := range g.orgs {
					if strings.HasPrefix(repoFullName, org+"/") {
						shouldInclude = true
						break
					}
				}
			}

			if shouldInclude {
				body := ""
				if includeBody && issue.Body != nil {
					body = *issue.Body
				}

				mappedIssue := Issue{
					Number:    *issue.Number,
					Title:     *issue.Title,
					Body:      body,
					State:     *issue.State,
					CreatedAt: issue.CreatedAt.Time,
					UpdatedAt: issue.UpdatedAt.Time,
					HTMLURL:   *issue.HTMLURL,
					Repository: struct {
						Name     string `json:"name"`
						FullName string `json:"full_name"`
					}{
						Name:     urlParts[4],
						FullName: repoFullName,
					},
					User: struct {
						Login string `json:"login"`
					}{
						Login: *issue.User.Login,
					},
				}
				allIssues = append(allIssues, mappedIssue)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allIssues, nil
}

func (g *Client) GetPullRequests(
	ctx context.Context,
	since, until time.Time,
	author string,
	includeBody bool,
) ([]PullRequest, error) {
	var allPRs []PullRequest

	// Search for pull requests created by the author
	query := fmt.Sprintf("author:%s created:%s..%s type:pr",
		author,
		since.Format("2006-01-02"),
		until.Format("2006-01-02"))

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		result, resp, err := g.client.Search.Issues(ctx, query, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to search pull requests: %w", err)
		}

		for _, issue := range result.Issues {
			// Extract repository info from the issue URL
			if issue.HTMLURL == nil {
				continue
			}

			// Parse repository from URL like https://github.com/owner/repo/pull/123
			urlParts := strings.Split(*issue.HTMLURL, "/")
			if len(urlParts) < 6 {
				continue
			}

			repoFullName := urlParts[3] + "/" + urlParts[4]

			// Check if the repository belongs to target organizations (if specified)
			shouldInclude := len(g.orgs) == 0 // Include all if no orgs specified
			if !shouldInclude {
				for _, org := range g.orgs {
					if strings.HasPrefix(repoFullName, org+"/") {
						shouldInclude = true
						break
					}
				}
			}

			if shouldInclude {
				var mergedAt *time.Time

				body := ""
				if includeBody && issue.Body != nil {
					body = *issue.Body
				}

				mappedPR := PullRequest{
					Number:    *issue.Number,
					Title:     *issue.Title,
					Body:      body,
					State:     *issue.State,
					CreatedAt: issue.CreatedAt.Time,
					UpdatedAt: issue.UpdatedAt.Time,
					MergedAt:  mergedAt,
					HTMLURL:   *issue.HTMLURL,
					Repository: struct {
						Name     string `json:"name"`
						FullName string `json:"full_name"`
					}{
						Name:     urlParts[4],
						FullName: repoFullName,
					},
					User: struct {
						Login string `json:"login"`
					}{
						Login: *issue.User.Login,
					},
				}
				allPRs = append(allPRs, mappedPR)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allPRs, nil
}

func (g *Client) GetReviewedPullRequests(
	ctx context.Context,
	since, until time.Time,
	reviewer string,
) ([]Review, error) {
	var allReviews []Review
	processedPRs := 0
	maxPRsToProcess := 100 // Limit to avoid timeout

	// Search for pull requests reviewed by the user
	query := fmt.Sprintf("reviewed-by:%s type:pr", reviewer)

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 50},
		Sort:        "updated",
		Order:       "desc",
	}

	for {
		result, resp, err := g.client.Search.Issues(ctx, query, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to search reviewed PRs: %w", err)
		}

		for _, issue := range result.Issues {
			// Limit processing to avoid timeout
			if processedPRs >= maxPRsToProcess {
				return allReviews, nil
			}

			// Extract repository info from the issue URL
			if issue.HTMLURL == nil {
				continue
			}

			// Parse repository from URL like https://github.com/owner/repo/pull/123
			urlParts := strings.Split(*issue.HTMLURL, "/")
			if len(urlParts) < 6 {
				continue
			}

			repoFullName := urlParts[3] + "/" + urlParts[4]

			// Check if the repository belongs to target organizations (if specified)
			shouldInclude := len(g.orgs) == 0 // Include all if no orgs specified
			if !shouldInclude {
				for _, org := range g.orgs {
					if strings.HasPrefix(repoFullName, org+"/") {
						shouldInclude = true
						break
					}
				}
			}
			if !shouldInclude {
				continue
			}

			processedPRs++

			// Get reviews for this PR
			parts := strings.Split(repoFullName, "/")
			if len(parts) == 2 {
				reviews, _, err := g.client.PullRequests.ListReviews(ctx, parts[0], parts[1], *issue.Number, nil)
				if err != nil {
					continue
				}

				// Also get review comments for this PR
				reviewComments, _, err := g.client.PullRequests.ListComments(
					ctx,
					parts[0],
					parts[1],
					*issue.Number,
					nil,
				)
				if err == nil {
					// Process review comments from the specified author
					for _, comment := range reviewComments {
						if comment.User != nil && comment.User.Login != nil &&
							*comment.User.Login == reviewer &&
							comment.CreatedAt != nil &&
							!comment.CreatedAt.Time.Before(since) &&
							!comment.CreatedAt.Time.After(until) {

							body := ""
							if comment.Body != nil {
								body = *comment.Body
							}

							prTitle := ""
							if issue.Title != nil {
								prTitle = *issue.Title
							}

							// Create a review entry for the comment
							mappedReview := Review{
								ID:          int(*comment.ID),
								State:       "COMMENTED",
								Body:        body,
								SubmittedAt: comment.CreatedAt.Time,
								HTMLURL:     *comment.HTMLURL,
								User: struct {
									Login string `json:"login"`
								}{
									Login: *comment.User.Login,
								},
								PullRequestURL:    *issue.HTMLURL,
								PullRequestTitle:  prTitle,
								PullRequestNumber: *issue.Number,
								Repository:        repoFullName,
							}
							allReviews = append(allReviews, mappedReview)
						}
					}
				}

				for _, review := range reviews {
					if review.User != nil && review.User.Login != nil &&
						*review.User.Login == reviewer &&
						review.SubmittedAt != nil &&
						!review.SubmittedAt.Time.Before(since) &&
						!review.SubmittedAt.Time.After(until) {

						body := ""
						if review.Body != nil {
							body = *review.Body
						}

						prTitle := ""
						if issue.Title != nil {
							prTitle = *issue.Title
						}

						mappedReview := Review{
							ID:          int(*review.ID),
							State:       *review.State,
							Body:        body,
							SubmittedAt: review.SubmittedAt.Time,
							HTMLURL:     *review.HTMLURL,
							User: struct {
								Login string `json:"login"`
							}{
								Login: *review.User.Login,
							},
							PullRequestURL:    *issue.HTMLURL,
							PullRequestTitle:  prTitle,
							PullRequestNumber: *issue.Number,
							Repository:        repoFullName,
						}
						allReviews = append(allReviews, mappedReview)
					}
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allReviews, nil
}
