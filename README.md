# ghdump

A tool to fetch GitHub activity for a specified period and generate Markdown reports.

## Features

- Fetch Issues within a specified period from specified organizations
- Fetch Pull Requests within a specified period from specified organizations
- Fetch reviewed Pull Requests within a specified period from specified organizations
- Generate reports in Markdown format

## Prerequisites

- Go
- GITHUB_TOKEN environment variable or GitHub CLI (`gh` command) authentication

## Installation

```bash
make build
```

## Usage

```bash
# Using GitHub CLI authentication token
GITHUB_TOKEN=$(gh auth token) ./ghdump -since 2025-01-01 -until 2025-06-30 -author username [-output report.md]

# Or set environment variable beforehand
export GITHUB_TOKEN=$(gh auth token)
./ghdump -since 2025-01-01 -until 2025-06-30 -author username [-output report.md]
```

### Options

- `-since`: Start date (YYYY-MM-DD format)
- `-until`: End date (YYYY-MM-DD format)
- `-author`: Target GitHub username
- `-output`: Output file path (outputs to stdout if not specified)
- `-body`: Include issue/PR descriptions (default: true)
- `-orgs`: Target GitHub organizations (comma-separated, targets all organizations if omitted)

### Examples

```bash
# Display report to stdout
GITHUB_TOKEN=$(gh auth token) ./ghdump -since 2024-01-01 -until 2024-03-31 -author username

# Save report to file
GITHUB_TOKEN=$(gh auth token) ./ghdump -since 2024-01-01 -until 2024-03-31 -author username -output quarterly-report.md

# Generate concise report without descriptions
GITHUB_TOKEN=$(gh auth token) ./ghdump -since 2024-01-01 -until 2024-03-31 -author username -body=false -output summary-report.md

# Target specific organizations only
GITHUB_TOKEN=$(gh auth token) ./ghdump -since 2024-01-01 -until 2024-03-31 -author username -orgs "myorg,anotherorg" -output specific-orgs-report.md
```

## Report Format

The generated report includes:

- Summary (period, author, counts for each section)
- Issues list (status, title, description, link, creation date, repository name)
- Pull Requests list (status, title, description, link, creation date, merge date, repository name)
- Reviews list (status, PR title, comments, review date, repository name)

## Notes

- GitHub API authentication is required (via GITHUB_TOKEN environment variable or GitHub CLI)
- Large data fetches may take time due to API rate limits
- Access permissions are required for repositories in specified organizations
- Review fetching is limited to a maximum of 100 PRs
