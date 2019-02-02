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

// NewListIssuesCommand will add a `github list-issues` command
func NewListIssuesCommand(ctx context.Context, client *github.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-issues",
		Short: "List all the Issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			return DisplayIssues(ctx, client, os.Stdout)
		},
	}

	return cmd
}

// DisplayIssues will display a list of issues
func DisplayIssues(ctx context.Context, client *github.Client, w io.Writer) error {
	githubOwner := viper.GetString("GITHUB_OWNER")

	allRepos, err := GetRepos(ctx, client, githubOwner)
	if err != nil {
		return err
	}

	rows := [][]string{}

	type DisplayIssue struct {
		RepoName string
		Title    string
	}
	issues := make(chan DisplayIssue)
	var wg sync.WaitGroup
	wg.Add(len(allRepos))

	go func() {
		wg.Wait()
		close(issues)
	}()

	for _, repo := range allRepos {
		go func(repo []string) {
			defer wg.Done()
			opt := &github.IssueListByRepoOptions{
				State: "open",
			}

			repoIssues, _, err := client.Issues.ListByRepo(ctx, repo[0], repo[1], opt)
			if err != nil {
				fmt.Fprintf(w, "unable to get issues for %s: %v", repo[1], err)
			}

			for _, repoIssue := range repoIssues {
				issue := DisplayIssue{
					RepoName: fmt.Sprintf("%s/%s", repo[0], repo[1]),
					Title:    fmt.Sprintf("#%v - %s", repoIssue.GetNumber(), repoIssue.GetTitle()),
				}
				issues <- issue
			}
		}(repo)
	}

	for result := range issues {
		rows = append(rows, []string{result.RepoName, result.Title})
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
