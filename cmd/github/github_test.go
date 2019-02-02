package github_test

import (
	"fmt"
	"testing"

	"github.com/benmatselby/frost/cmd/github"
	"github.com/benmatselby/frost/test"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestNewGitHubCommand(t *testing.T) {
	cmd := github.NewGitHubCommand()

	expected := &cobra.Command{
		Use:   "github",
		Short: "GitHub related commands",
	}

	test.Command(t, cmd, expected)
}

func TestShowRepoRules(t *testing.T) {
	tt := []struct {
		name      string
		org       string
		repoName  string
		useConfig bool
		expected  bool
	}{
		{name: "standard org repo check", org: "apple", repoName: "mac", useConfig: true, expected: true},
		{name: "a microsoft repo", org: "microsoft", repoName: "vscode", useConfig: true, expected: true},
		{name: "a repo not defined", org: "google", repoName: "mail", useConfig: true, expected: false},
		{name: "a fork of golang", org: "bobdylan", repoName: "golang", useConfig: true, expected: true},
		{name: "any repo as config not used", org: "microsoft", repoName: "office", useConfig: false, expected: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			if tc.useConfig {
				viper.Set("github.repos", []string{
					"apple/mac",
					"microsoft/*",
					"*/golang",
				})
			} else {
				viper.Set("github.repos", []string{})
			}

			result := github.ShowRepo(tc.org, tc.repoName)

			if result != tc.expected {
				t.Fatalf("expected %v for %s; got %v", tc.expected, fmt.Sprintf("%s/%s", tc.org, tc.repoName), result)
			}
		})
	}
}
