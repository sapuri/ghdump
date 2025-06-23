package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sapuri/ghdump/internal/github"
	"github.com/sapuri/ghdump/internal/reporter"
)

func main() {
	var (
		since       = flag.String("since", "", "Start date (YYYY-MM-DD)")
		until       = flag.String("until", "", "End date (YYYY-MM-DD)")
		author      = flag.String("author", "", "GitHub username to filter by")
		output      = flag.String("output", "", "Output file path (optional)")
		includeBody = flag.Bool("body", true, "Include issue/PR descriptions")
		orgs        = flag.String("orgs", "", "Comma-separated list of GitHub organizations (optional: if not specified, searches all organizations)")
	)
	flag.Parse()

	if *since == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Error: -since is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *until == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Error: -until is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *author == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Error: -author is required\n")
		flag.Usage()
		os.Exit(1)
	}

	sinceTime, err := time.Parse("2006-01-02", *since)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing since date: %v\n", err)
		os.Exit(1)
	}

	untilTime, err := time.Parse("2006-01-02", *until)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing until date: %v\n", err)
		os.Exit(1)
	}

	var orgList []string
	if *orgs != "" {
		orgList = strings.Split(*orgs, ",")
		for i, org := range orgList {
			orgList[i] = strings.TrimSpace(org)
		}
	}

	client, err := github.NewClient(orgList)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error creating GitHub client: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Display configuration
	_, _ = fmt.Fprintf(os.Stderr, "Configuration:\n")
	_, _ = fmt.Fprintf(os.Stderr, "  Period: %s - %s\n", *since, *until)
	_, _ = fmt.Fprintf(os.Stderr, "  Author: %s\n", *author)
	_, _ = fmt.Fprintf(os.Stderr, "  Include body: %t\n", *includeBody)
	if *orgs != "" {
		_, _ = fmt.Fprintf(os.Stderr, "  Organizations: %s\n", *orgs)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "  Organizations: all\n")
	}
	if *output != "" {
		_, _ = fmt.Fprintf(os.Stderr, "  Output: %s\n", *output)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "  Output: stdout\n")
	}
	_, _ = fmt.Fprintf(os.Stderr, "\n")

	_, _ = fmt.Fprintf(os.Stderr, "Fetching issues...\n")
	issues, err := client.GetIssues(ctx, sinceTime, untilTime, *author, *includeBody)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error fetching issues: %v\n", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Fetching pull requests...\n")
	prs, err := client.GetPullRequests(ctx, sinceTime, untilTime, *author, *includeBody)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error fetching pull requests: %v\n", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Fetching reviews...\n")
	reviews, err := client.GetReviewedPullRequests(ctx, sinceTime, untilTime, *author)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error fetching reviews: %v\n", err)
		os.Exit(1)
	}

	report := reporter.New(*includeBody).GenerateMarkdownReport(
		issues,
		prs,
		reviews,
		sinceTime,
		untilTime,
		*author,
	)

	if *output != "" {
		if err := os.WriteFile(*output, []byte(report), 0644); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		_, _ = fmt.Fprintf(os.Stderr, "Report written to %s\n", *output)
	} else {
		fmt.Print(report)
	}
}
