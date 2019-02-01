package github

import (
	"context"
	"fmt"
	"strings"

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
		NewListIssuesCommand(ctx, client),
	)

	return cmd
}

// GetRepos will find all repos the github owner can "see"
func GetRepos(ctx context.Context, client *github.Client, githubOwner string) ([][]string, error) {
	orgs, _, err := client.Organizations.List(context.Background(), githubOwner, nil)
	if err != nil {
		return nil, err
	}

	var allRepos [][]string
	for _, org := range orgs {
		opt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}
		repos, _, err := client.Repositories.ListByOrg(ctx, org.GetLogin(), opt)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			if !ShowRepo(org.GetLogin(), repo.GetName()) {
				continue
			}
			allRepos = append(allRepos, []string{org.GetLogin(), repo.GetName()})
		}
	}

	opt := &github.RepositoryListOptions{}
	repos, _, err := client.Repositories.List(ctx, githubOwner, opt)
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if repo.GetFork() {
			continue
		}
		if !ShowRepo(githubOwner, repo.GetName()) {
			continue
		}
		allRepos = append(allRepos, []string{githubOwner, repo.GetName()})
	}

	return allRepos, nil
}

// ShowRepo is going to determine if we care enough to show the detail
func ShowRepo(org, name string) bool {
	watchRepos := viper.GetStringSlice("github.repos")

	if len(watchRepos) == 0 {
		return true
	}

	show := false
	for _, i := range watchRepos {
		s := strings.Split(i, "/")

		// If we want to watch everything for a given org
		if s[1] == "*" && org == s[0] {
			show = true
			break
		}

		// If we want to watch everything for a given repo (including forks)
		if s[0] == "*" && name == s[1] {
			show = true
			break
		}

		// Otherwise we want an exact match
		if i == fmt.Sprintf("%s/%s", org, name) {
			show = true
			break
		}
	}

	return show
}
