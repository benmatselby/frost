package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	travis "github.com/Ableton/go-travis"
	"github.com/urfave/cli"
)

func travisListBuildOverview(c *cli.Context) {
	if len(travisOwner) <= 0 {
		fmt.Fprintf(os.Stderr, "os.Env TRAVIS_CI_OWNER not defined")
		os.Exit(2)
	}

	if len(travisToken) <= 0 {
		fmt.Fprintf(os.Stderr, "os.Env TRAVIS_CI_TOKEN not defined")
		os.Exit(2)
	}

	filterBranch := c.String("branch")

	client := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, travisToken)
	opt := &travis.RepositoryListOptions{Member: travisOwner, Active: true}
	repos, _, err := client.Repositories.Find(opt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get the overview: %v", err)
		os.Exit(2)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "", "Name", "Branch", "Finished")
	for _, repo := range repos {
		// Trying to remove the items that are not really running in Travis CI
		// Assume there is a better way to do this?
		if repo.LastBuildState == "" {
			continue
		}

		for _, branchName := range strings.Split(filterBranch, ",") {
			branch, _, err := client.Branches.GetFromSlug(repo.Slug, branchName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to get the branches: %v", err)
				os.Exit(2)
			}

			finish, err := time.Parse(time.RFC3339, branch.FinishedAt)
			finishAt := finish.Format(appDateTimeFormat)
			if err != nil {
				finishAt = branch.FinishedAt
			}

			result := ""
			if branch.State == "failed" {
				result = appFailure
			} else if branch.State == "started" {
				result = appProgress
			} else {
				result = appSuccess
			}

			fmt.Fprintf(w, "%s \t%s\t%s\t%s\n", result, repo.Slug, branchName, finishAt)
		}
	}

	w.Flush()
}
