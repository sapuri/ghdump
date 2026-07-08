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
		"repo filter takes precedence over unrelated org filter": {
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
