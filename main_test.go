package main

import (
	"slices"
	"testing"
)

func TestSplitAndTrim(t *testing.T) {
	cases := map[string]struct {
		input string
		want  []string
	}{
		"empty string returns nil": {input: "", want: nil},
		"single value":             {input: "myorg", want: []string{"myorg"}},
		"multiple values trimmed":  {input: "myorg, another , third", want: []string{"myorg", "another", "third"}},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := splitAndTrim(tc.input); !slices.Equal(got, tc.want) {
				t.Errorf("splitAndTrim(%q) = %#v, want %#v", tc.input, got, tc.want)
			}
		})
	}
}
