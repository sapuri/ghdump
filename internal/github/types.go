package github

import "time"

// Issue represents a GitHub issue.
type Issue struct {
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	State      string    `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	HTMLURL    string    `json:"html_url"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
	User struct {
		Login string `json:"login"`
	} `json:"user"`
}

// PullRequest represents a GitHub pull request.
type PullRequest struct {
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	State      string     `json:"state"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	MergedAt   *time.Time `json:"merged_at"`
	HTMLURL    string     `json:"html_url"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
	User struct {
		Login string `json:"login"`
	} `json:"user"`
}

// Review represents a GitHub pull request review.
type Review struct {
	ID          int       `json:"id"`
	State       string    `json:"state"`
	Body        string    `json:"body"`
	SubmittedAt time.Time `json:"submitted_at"`
	HTMLURL     string    `json:"html_url"`
	User        struct {
		Login string `json:"login"`
	} `json:"user"`
	PullRequestURL    string `json:"pull_request_url"`
	PullRequestTitle  string `json:"pull_request_title"`
	PullRequestNumber int    `json:"pull_request_number"`
	Repository        string `json:"repository"`
}
