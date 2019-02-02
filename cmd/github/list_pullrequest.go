package github

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"text/tabwriter"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewPullRequestCommand will add a `github list-pr` command
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

	allRepos, err := GetRepos(ctx, client, githubOwner)
	if err != nil {
		return err
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

	tr := tabwriter.NewWriter(w, 0, 1, 1, ' ', 0)
	fmt.Fprintf(tr, "%s\t%s\n", "Repo", "Title")
	for _, row := range rows {
		fmt.Fprintf(tr, "%s\t%s\n", row[0], row[1])
	}
	tr.Flush()

	return nil
}
