package build

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	travis "github.com/Ableton/go-travis"
	"github.com/benmatselby/frost/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewTravisBuildListCommand will add a `travis build list` command which
// is responsible for showing a list of build
func NewTravisBuildListCommand(client *travis.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the builds",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ListBuilds(client, os.Stdout)
		},
	}

	return cmd
}

// ListBuilds will display a list of builds
func ListBuilds(client *travis.Client, w io.Writer) error {
	filterBranch := "master"

	opt := &travis.RepositoryListOptions{Member: viper.GetString("TRAVIS_CI_OWNER"), Active: true}
	repos, _, err := client.Repositories.Find(opt)
	if err != nil {
		return err
	}

	tr := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(tr, "%s\t%s\t%s\t%s\n", "", "Name", "Branch", "Finished")
	for _, repo := range repos {
		// Trying to remove the items that are not really running in Travis CI
		// Assume there is a better way to do this?
		if repo.LastBuildState == "" {
			continue
		}

		for _, branchName := range strings.Split(filterBranch, ",") {
			branch, _, err := client.Branches.GetFromSlug(repo.Slug, branchName)
			if err != nil {
				return err
			}

			finish, err := time.Parse(time.RFC3339, branch.FinishedAt)
			finishAt := finish.Format(ui.AppDateTimeFormat)
			if err != nil {
				finishAt = branch.FinishedAt
			}

			result := ""
			if branch.State == "failed" {
				result = ui.AppFailure
			} else if branch.State == "started" {
				result = ui.AppProgress
			} else {
				result = ui.AppSuccess
			}

			fmt.Fprintf(tr, "%s \t%s\t%s\t%s\n", result, repo.Slug, branchName, finishAt)
		}
	}

	tr.Flush()

	return nil
}
