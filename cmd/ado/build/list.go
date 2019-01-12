package build

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/frost/ui"
	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListOptions provides the flags for the `vet` command
type ListOptions struct {
	Branches string
	flags    []string
}

// NewAdoBuildListCommand will add a `ado build list` command which is responsible
// for showing a list of build
func NewAdoBuildListCommand(client *azuredevops.Client) *cobra.Command {
	var opts ListOptions
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the builds",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.flags = args
			return ListBuilds(client, opts, os.Stdout)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Branches, "branches", "master", "Which branches should be displayed")

	return cmd
}

// ListBuilds will display a list of builds
func ListBuilds(client *azuredevops.Client, opts ListOptions, w io.Writer) error {
	buildDefOpts := azuredevops.BuildDefinitionsListOptions{Path: "\\" + viper.GetString("AZURE_DEVOPS_TEAM")}
	definitions, err := client.BuildDefinitions.List(&buildDefOpts)
	if err != nil {
		return err
	}

	results := make(chan azuredevops.Build)
	var wg sync.WaitGroup
	wg.Add(len(definitions))

	go func() {
		wg.Wait()
		close(results)
	}()

	for _, definition := range definitions {
		go func(definition azuredevops.BuildDefinition) {
			defer wg.Done()

			for _, branchName := range strings.Split(opts.Branches, ",") {
				builds, err := getBuildsForBranch(client, definition.ID, branchName)
				if err != nil {
					fmt.Fprintf(w, "unable to get builds for definition %s: %v", definition.Name, err)
				}
				if len(builds) > 0 {
					results <- builds[0]
				}
			}
		}(definition)
	}

	var builds []azuredevops.Build
	for result := range results {
		builds = append(builds, result)
	}

	sort.Slice(builds, func(i, j int) bool { return builds[i].Definition.Name < builds[j].Definition.Name })

	// renderAzureDevOpsBuilds(builds, len(builds), ".*")
	tr := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "", "Name", "Branch", "Build", "Finished")
	for index := 0; index < len(builds); index++ {
		build := builds[index]
		name := build.Definition.Name
		result := build.Result
		status := build.Status
		buildNo := build.BuildNumber
		branch := build.Branch

		// Deal with date formatting for the finish time
		finish, err := time.Parse(time.RFC3339, builds[index].FinishTime)
		finishAt := finish.Format(ui.AppDateTimeFormat)
		if err != nil {
			finishAt = builds[index].FinishTime
		}

		// Filter on branches
		matched, _ := regexp.MatchString(".*"+opts.Branches+".*", branch)
		if matched == false {
			continue
		}

		if status == "inProgress" {
			result = ui.AppProgress
		} else if status == "notStarted" {
			result = ui.AppPending
		} else {
			if result == "failed" {
				result = ui.AppFailure
			} else {
				result = ui.AppSuccess
			}
		}

		fmt.Fprintf(tr, "%s \t%s\t%s\t%s\t%s\n", result, name, branch, buildNo, finishAt)
	}

	tr.Flush()

	return nil
}

func getBuildsForBranch(client *azuredevops.Client, defID int, branchName string) ([]azuredevops.Build, error) {
	buildOpts := azuredevops.BuildsListOptions{Definitions: strconv.Itoa(defID), Branch: "refs/heads/" + branchName, Count: 1}
	build, err := client.Builds.List(&buildOpts)
	return build, err
}
