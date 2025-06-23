# ghdump

A tool to fetch GitHub activity for a specified period and generate Markdown reports.

## Features

- Fetch Issues within a specified period from specified organizations
- Fetch Pull Requests within a specified period from specified organizations
- Fetch reviewed Pull Requests within a specified period from specified organizations
- Generate reports in Markdown format

## Prerequisites

- Go
- `GITHUB_TOKEN` environment variable or GitHub CLI (`gh` command) authentication

## Installation

### Using go install

```bash
go install github.com/sapuri/ghdump@latest
```

### Build from source

```bash
git clone https://github.com/sapuri/ghdump.git
cd ghdump
make build
```

## Usage

```bash
# Using GitHub CLI authentication token
GITHUB_TOKEN=$(gh auth token) ghdump -since 2025-01-01 -until 2025-06-30 -author username [-output report.md]

# Or set environment variable beforehand
export GITHUB_TOKEN=$(gh auth token)
ghdump -since 2025-01-01 -until 2025-06-30 -author username [-output report.md]
```

### Options

```
$ ghdump -h
Usage of ghdump:
  -author string
    	GitHub username to filter by
  -body
    	Include issue/PR descriptions (default true)
  -orgs string
    	Comma-separated list of GitHub organizations (optional: if not specified, searches all organizations)
  -output string
    	Output file path (optional)
  -since string
    	Start date (YYYY-MM-DD)
  -until string
    	End date (YYYY-MM-DD)
```

### Examples

```bash
# Display report to stdout
GITHUB_TOKEN=$(gh auth token) ghdump -since 2025-01-01 -until 2025-06-30 -author username

# Save report to file
GITHUB_TOKEN=$(gh auth token) ghdump -since 2025-01-01 -until 2025-06-30 -author username -output quarterly-report.md

# Generate concise report without descriptions
GITHUB_TOKEN=$(gh auth token) ghdump -since 2025-01-01 -until 2025-06-30 -author username -body=false -output summary-report.md

# Target specific organizations only
GITHUB_TOKEN=$(gh auth token) ghdump -since 2025-01-01 -until 2025-06-30 -author username -orgs "myorg,anotherorg" -output specific-orgs-report.md
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
