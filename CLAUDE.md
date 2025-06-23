# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

- **Build**: `make build` - Builds the binary to `./ghdump`
- **Run**: `./ghdump -since YYYY-MM-DD -until YYYY-MM-DD -author username` - Run the built binary
- **Test**: No test suite is configured in this repository

## Architecture Overview

This is a Go CLI application that fetches GitHub activity data and generates Markdown reports. The application has three main packages:

### Core Components

1. **main.go**: CLI entry point that parses flags and orchestrates the data fetching and report generation
2. **internal/github**: GitHub API client package that handles authentication and data fetching
3. **internal/reporter**: Report generation package that formats data into Markdown

### Data Flow

1. CLI validates required flags (-since, -until, -author) and parses date ranges
2. GitHub client authenticates using GITHUB_TOKEN environment variable
3. Three parallel API calls fetch:
   - Issues created by the author
   - Pull requests created by the author  
   - Pull request reviews performed by the author
4. Reporter formats all data into a structured Markdown report

### Key Implementation Details

- Uses GitHub Search API for issues and PRs to filter by author and date range
- Organization filtering is applied client-side by parsing repository URLs
- Review fetching is limited to 100 PRs maximum to avoid API timeouts
- Both PR reviews and review comments are collected for comprehensive review data
- Report includes status indicators (🔴 open, ✅ closed/approved, 🟣 merged, 🔄 changes requested)

### Authentication

Requires GITHUB_TOKEN environment variable. Common pattern:
```bash
export GITHUB_TOKEN=$(gh auth token)
ghdump [flags]
```

### Dependencies

- `github.com/google/go-github/v72` - GitHub API client
- `golang.org/x/oauth2` - OAuth2 authentication
