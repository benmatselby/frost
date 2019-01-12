package github

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewPullRequestCommand will add a `github pullrequest list` command
func NewPullRequestCommand(ctx context.Context, client *github.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-pr",
		Short: "List all the Pull Requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			return List(ctx, client, os.Stdout)
		},
	}

	return cmd
}

// List will display a list of pull requests
func List(ctx context.Context, client *github.Client, w io.Writer) error {
	githubOwner := viper.GetString("GITHUB_OWNER")

	orgs, _, err := client.Organizations.List(context.Background(), githubOwner, nil)
	if err != nil {
		return err
	}

	var allRepos [][]string
	for _, org := range orgs {
		opt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}
		repos, _, err := client.Repositories.ListByOrg(ctx, org.GetLogin(), opt)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			if !showRepoPr(org.GetLogin(), repo.GetName()) {
				continue
			}
			allRepos = append(allRepos, []string{org.GetLogin(), repo.GetName()})
		}
	}

	opt := &github.RepositoryListOptions{}
	repos, _, err := client.Repositories.List(ctx, githubOwner, opt)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		if repo.GetFork() {
			continue
		}
		if !showRepoPr(githubOwner, repo.GetName()) {
			continue
		}
		allRepos = append(allRepos, []string{githubOwner, repo.GetName()})
	}

	rows := [][]string{}

	pullRequests := make(chan *github.PullRequest)
	var wg sync.WaitGroup
	wg.Add(len(allRepos))

	go func() {
		wg.Wait()
		close(pullRequests)
	}()

	for _, repo := range allRepos {
		go func(repo []string) {
			defer wg.Done()
			opt := &github.PullRequestListOptions{
				State: "open",
			}

			prs, _, err := client.PullRequests.List(ctx, repo[0], repo[1], opt)
			if err != nil {
				fmt.Fprintf(w, "unable to get pull requests for %s: %v", repo[1], err)
			}

			for _, pull := range prs {
				pullRequests <- pull
			}
		}(repo)
	}

	for result := range pullRequests {
		rows = append(rows, []string{result.GetHead().GetRepo().GetFullName(), fmt.Sprintf("#%v - %s", result.GetNumber(), result.GetTitle())})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	tr := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(tr, "%s\t%s\n", "Repo", "Title.")
	for _, row := range rows {
		fmt.Fprintf(tr, "%s\t%s\n", row[0], row[1])
	}
	tr.Flush()

	return nil
}

// showRepoPr is going to determine if we care enough to show the detail
func showRepoPr(org, name string) bool {
	watchRepos := viper.GetStringSlice("github.pull_request_repos")

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
