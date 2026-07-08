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
		"empty string returns nil":    {input: "", want: nil},
		"single value":                {input: "myorg", want: []string{"myorg"}},
		"multiple values trimmed":     {input: "myorg, another , third", want: []string{"myorg", "another", "third"}},
		"trailing comma drops empty":  {input: "myorg,", want: []string{"myorg"}},
		"whitespace-only returns nil": {input: "  , ", want: nil},
		"leading comma drops empty":   {input: ",myorg", want: []string{"myorg"}},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := splitAndTrim(tc.input); !slices.Equal(got, tc.want) {
				t.Errorf("splitAndTrim(%q) = %#v, want %#v", tc.input, got, tc.want)
			}
		})
	}
}
