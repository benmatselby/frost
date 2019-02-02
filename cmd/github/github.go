package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/spf13/cobra"
)

// NewGitHubCommand will add a `github` command which will add sub commands
func NewGitHubCommand() *cobra.Command {
	ctx, client := NewClient()

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

// NewClient returns a github client
func NewClient() (context.Context, *github.Client) {
	githubToken := viper.GetString("GITHUB_TOKEN")
	ctx := context.Background()

	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	httpClient := oauth2.NewClient(ctx, sts)
	client := github.NewClient(httpClient)

	return ctx, client
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
	repos := getConfigWatchRepos()

	if len(repos) == 0 {
		return true
	}

	slug := fmt.Sprintf("%s/%s", org, name)

	if repos["!"+slug] == true {
		return false
	}

	if repos[slug] == true {
		return true
	}

	if repos[org+"/*"] == true {
		return true
	}

	if repos["*/"+name] == true {
		return true
	}

	return false
}

var configWatchRepos = map[string]bool{}
var configWatchRepoSet = false

func getConfigWatchRepos() map[string]bool {
	watchRepos := viper.GetStringSlice("github.repos")

	if configWatchRepoSet == true {
		return configWatchRepos
	}

	for _, i := range watchRepos {
		configWatchRepos[i] = true
	}

	configWatchRepoSet = true
	return configWatchRepos
}
