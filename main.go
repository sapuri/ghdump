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
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		since       = flag.String("since", "", "Start date (YYYY-MM-DD)")
		until       = flag.String("until", "", "End date (YYYY-MM-DD)")
		author      = flag.String("author", "", "GitHub username to filter by")
		output      = flag.String("output", "", "Output file path (optional)")
		includeBody = flag.Bool("body", true, "Include issue/PR descriptions")
		orgs        = flag.String("orgs", "", "Comma-separated list of GitHub organizations (optional: if not specified, searches all organizations)")
		repos       = flag.String("repos", "", "Comma-separated list of GitHub repositories in owner/repo format (optional: if not specified, searches all repositories)")
	)
	flag.Parse()

	requiredFlags := []struct {
		name  string
		value string
	}{
		{"since", *since},
		{"until", *until},
		{"author", *author},
	}
	for _, f := range requiredFlags {
		if f.value == "" {
			flag.Usage()
			return fmt.Errorf("-%s is required", f.name)
		}
	}

	sinceTime, err := time.Parse("2006-01-02", *since)
	if err != nil {
		return fmt.Errorf("parsing since date: %w", err)
	}

	untilTime, err := time.Parse("2006-01-02", *until)
	if err != nil {
		return fmt.Errorf("parsing until date: %w", err)
	}

	orgList := splitAndTrim(*orgs)
	repoList := splitAndTrim(*repos)

	client, err := github.NewClient(orgList, repoList)
	if err != nil {
		return fmt.Errorf("creating GitHub client: %w", err)
	}

	printConfiguration(*since, *until, *author, *includeBody, *orgs, *repos, *output)

	ctx := context.Background()

	_, _ = fmt.Fprintf(os.Stderr, "Fetching issues...\n")
	issues, err := client.GetIssues(ctx, sinceTime, untilTime, *author, *includeBody)
	if err != nil {
		return fmt.Errorf("fetching issues: %w", err)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Fetching pull requests...\n")
	prs, err := client.GetPullRequests(ctx, sinceTime, untilTime, *author, *includeBody)
	if err != nil {
		return fmt.Errorf("fetching pull requests: %w", err)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Fetching reviews...\n")
	reviews, err := client.GetReviewedPullRequests(ctx, sinceTime, untilTime, *author)
	if err != nil {
		return fmt.Errorf("fetching reviews: %w", err)
	}

	report := reporter.New(*includeBody).GenerateMarkdownReport(
		issues,
		prs,
		reviews,
		sinceTime,
		untilTime,
		*author,
	)

	if *output == "" {
		fmt.Print(report)
		return nil
	}

	if err := os.WriteFile(*output, []byte(report), 0644); err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Report written to %s\n", *output)

	return nil
}

// printConfiguration prints the resolved run configuration to stderr.
func printConfiguration(since, until, author string, includeBody bool, orgs, repos, output string) {
	_, _ = fmt.Fprintf(os.Stderr, "Configuration:\n")
	_, _ = fmt.Fprintf(os.Stderr, "  Period: %s - %s\n", since, until)
	_, _ = fmt.Fprintf(os.Stderr, "  Author: %s\n", author)
	_, _ = fmt.Fprintf(os.Stderr, "  Include body: %t\n", includeBody)
	printFilterLine("Organizations", orgs, "all")
	printFilterLine("Repositories", repos, "all")
	printFilterLine("Output", output, "stdout")
	_, _ = fmt.Fprintf(os.Stderr, "\n")
}

// printFilterLine prints a single "label: value" configuration line, falling
// back to defaultValue when value is empty, e.g. printFilterLine("Output",
// "", "stdout") prints "  Output: stdout".
func printFilterLine(label, value, defaultValue string) {
	if value == "" {
		value = defaultValue
	}
	_, _ = fmt.Fprintf(os.Stderr, "  %s: %s\n", label, value)
}

// splitAndTrim splits a comma-separated flag value into trimmed, non-empty
// elements, returning nil if none remain (e.g. for an empty string, a
// whitespace-only value, or a trailing comma).
func splitAndTrim(s string) []string {
	var result []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			result = append(result, p)
		}
	}
	return result
}
