package github

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// NewGitHubCommand will add a `github` command which will add sub commands
func NewGitHubCommand() *cobra.Command {
	githubToken := viper.GetString("GITHUB_TOKEN")
	ctx := context.Background()

	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	httpClient := oauth2.NewClient(ctx, sts)
	client := github.NewClient(httpClient)

	cmd := &cobra.Command{
		Use:   "github",
		Short: "GitHub related commands",
	}

	cmd.AddCommand(
		NewPullRequestCommand(ctx, client),
	)

	return cmd
}
