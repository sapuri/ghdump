package github

import "testing"

func TestShouldIncludeRepo(t *testing.T) {
	cases := map[string]struct {
		orgs         []string
		repos        []string
		repoFullName string
		want         bool
	}{
		"no filters includes everything": {
			repoFullName: "myorg/myrepo",
			want:         true,
		},
		"org match": {
			orgs:         []string{"myorg"},
			repoFullName: "myorg/myrepo",
			want:         true,
		},
		"org mismatch": {
			orgs:         []string{"otherorg"},
			repoFullName: "myorg/myrepo",
			want:         false,
		},
		"org match is case-insensitive": {
			orgs:         []string{"MyOrg"},
			repoFullName: "myorg/myrepo",
			want:         true,
		},
		"org prefix does not falsely match a similarly named org": {
			orgs:         []string{"my"},
			repoFullName: "myorg/myrepo",
			want:         false,
		},
		"repo exact match": {
			repos:        []string{"myorg/myrepo"},
			repoFullName: "myorg/myrepo",
			want:         true,
		},
		"repo mismatch": {
			repos:        []string{"myorg/other"},
			repoFullName: "myorg/myrepo",
			want:         false,
		},
		"repo match is case-insensitive": {
			repos:        []string{"MyOrg/MyRepo"},
			repoFullName: "myorg/myrepo",
			want:         true,
		},
		"org mismatch excludes despite a matching repo filter": {
			orgs:         []string{"otherorg"},
			repos:        []string{"myorg/myrepo"},
			repoFullName: "myorg/myrepo",
			want:         false,
		},
		"repo and org filters combine with AND": {
			orgs:         []string{"myorg"},
			repos:        []string{"myorg/myrepo"},
			repoFullName: "myorg/myrepo",
			want:         true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := &Client{orgs: tc.orgs, repos: tc.repos}
			if got := c.shouldIncludeRepo(tc.repoFullName); got != tc.want {
				t.Errorf("shouldIncludeRepo(%q) = %v, want %v", tc.repoFullName, got, tc.want)
			}
		})
	}
}

func TestIsValidRepoFullName(t *testing.T) {
	cases := map[string]struct {
		repo string
		want bool
	}{
		"valid owner/repo":  {repo: "myorg/myrepo", want: true},
		"missing owner":     {repo: "myrepo", want: false},
		"empty owner":       {repo: "/myrepo", want: false},
		"empty repo":        {repo: "myorg/", want: false},
		"too many segments": {repo: "myorg/myrepo/extra", want: false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := isValidRepoFullName(tc.repo); got != tc.want {
				t.Errorf("isValidRepoFullName(%q) = %v, want %v", tc.repo, got, tc.want)
			}
		})
	}
}

func TestRepoFullNameFromURL(t *testing.T) {
	url := "https://github.com/myorg/myrepo/issues/123"

	cases := map[string]struct {
		htmlURL  *string
		wantName string
		wantOK   bool
	}{
		"valid issue URL": {htmlURL: &url, wantName: "myorg/myrepo", wantOK: true},
		"nil URL":         {htmlURL: nil, wantName: "", wantOK: false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gotName, gotOK := repoFullNameFromURL(tc.htmlURL)
			if gotName != tc.wantName || gotOK != tc.wantOK {
				t.Errorf("repoFullNameFromURL(%v) = (%q, %v), want (%q, %v)", tc.htmlURL, gotName, gotOK, tc.wantName, tc.wantOK)
			}
		})
	}
}
