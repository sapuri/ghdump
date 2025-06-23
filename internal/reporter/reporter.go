package reporter

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sapuri/ghdump/internal/github"
)

// Reporter generates activity reports.
type Reporter struct {
	includeBody bool
}

// New creates a new Reporter.
func New(includeBody bool) *Reporter {
	return &Reporter{
		includeBody: includeBody,
	}
}

// GenerateMarkdownReport creates a Markdown report from GitHub activity data.
func (r *Reporter) GenerateMarkdownReport(
	issues []github.Issue,
	prs []github.PullRequest,
	reviews []github.Review,
	since, until time.Time,
	author string,
) string {
	var sb strings.Builder

	periodStr := fmt.Sprintf("%s - %s", since.Format("2006-01-02"), until.Format("2006-01-02"))

	sb.WriteString("# GitHub Activity Report\n\n")
	sb.WriteString(fmt.Sprintf("**Period:** %s\n\n", periodStr))
	sb.WriteString(fmt.Sprintf("**Author:** %s\n\n", author))

	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Issues:** %d\n", len(issues)))
	sb.WriteString(fmt.Sprintf("- **Pull Requests:** %d\n", len(prs)))
	sb.WriteString(fmt.Sprintf("- **Reviews:** %d\n\n", len(reviews)))

	if len(issues) > 0 {
		r.writeIssues(&sb, issues)
	}

	if len(prs) > 0 {
		r.writePullRequests(&sb, prs)
	}

	if len(reviews) > 0 {
		r.writeReviews(&sb, reviews)
	}

	return sb.String()
}

func (r *Reporter) writeIssues(sb *strings.Builder, issues []github.Issue) {
	// Sort issues by created date (oldest first)
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].CreatedAt.Before(issues[j].CreatedAt)
	})

	sb.WriteString("## Issues\n\n")
	for _, issue := range issues {
		status := "🔴"
		if issue.State == "closed" {
			status = "✅"
		}
		sb.WriteString(fmt.Sprintf(
			"### %s [#%d %s](%s)\n",
			status,
			issue.Number,
			issue.Title,
			issue.HTMLURL,
		))
		sb.WriteString(fmt.Sprintf("- **Repository:** %s\n", issue.Repository.FullName))
		sb.WriteString(fmt.Sprintf("- **Created:** %s\n", issue.CreatedAt.Format("2006-01-02")))

		if r.includeBody && issue.Body != "" {
			sb.WriteString("- **Description:**\n")
			r.writeIndentedText(sb, issue.Body)
		}
		sb.WriteString("\n")
	}
}

func (r *Reporter) writePullRequests(sb *strings.Builder, prs []github.PullRequest) {
	// Sort PRs by created date (oldest first)
	sort.Slice(prs, func(i, j int) bool {
		return prs[i].CreatedAt.Before(prs[j].CreatedAt)
	})

	sb.WriteString("## Pull Requests\n\n")
	for _, pr := range prs {
		status := "🔴"
		if pr.State == "merged" {
			status = "🟣"
		} else if pr.State == "closed" {
			status = "✅"
		}

		sb.WriteString(fmt.Sprintf("### %s [#%d %s](%s)\n",
			status, pr.Number, pr.Title, pr.HTMLURL))
		sb.WriteString(fmt.Sprintf("- **Repository:** %s\n", pr.Repository.FullName))
		sb.WriteString(fmt.Sprintf("- **Created:** %s\n", pr.CreatedAt.Format("2006-01-02")))

		if pr.MergedAt != nil {
			sb.WriteString(fmt.Sprintf("- **Merged:** %s\n", pr.MergedAt.Format("2006-01-02")))
		}

		if r.includeBody && pr.Body != "" {
			sb.WriteString("- **Description:**\n")
			r.writeIndentedText(sb, pr.Body)
		}
		sb.WriteString("\n")
	}
}

func (r *Reporter) writeReviews(sb *strings.Builder, reviews []github.Review) {
	// Sort reviews by submitted date (oldest first)
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].SubmittedAt.Before(reviews[j].SubmittedAt)
	})

	sb.WriteString("## Reviews\n\n")
	for _, review := range reviews {
		status := "💬"
		if review.State == "APPROVED" {
			status = "✅"
		} else if review.State == "CHANGES_REQUESTED" {
			status = "🔄"
		}

		sb.WriteString(fmt.Sprintf("### %s [#%d %s](%s)\n",
			status, review.PullRequestNumber, review.PullRequestTitle, review.PullRequestURL))
		sb.WriteString(fmt.Sprintf("- **Repository:** %s\n", review.Repository))
		sb.WriteString(fmt.Sprintf("- **Reviewed:** %s\n", review.SubmittedAt.Format("2006-01-02")))

		if review.Body != "" {
			sb.WriteString("- **Comment:**\n")
			r.writeIndentedText(sb, review.Body)
		}
		sb.WriteString("\n")
	}
}

func (r *Reporter) writeIndentedText(sb *strings.Builder, text string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			sb.WriteString(fmt.Sprintf("  %s\n", line))
		}
	}
}
