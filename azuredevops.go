package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/urfave/cli"
)

func azureListBuildOverview(c *cli.Context) {
	if len(azureDevOpsAccount) <= 0 {
		fmt.Fprintf(os.Stderr, "os.Env AZURE_DEVOPS_ACCOUNT not defined")
		os.Exit(2)
	}

	if len(azureDevOpsProject) <= 0 {
		fmt.Fprintf(os.Stderr, "os.Env AZURE_DEVOPS_PROJECT not defined")
		os.Exit(2)
	}

	if len(azureDevOpsToken) <= 0 {
		fmt.Fprintf(os.Stderr, "os.Env AZURE_DEVOPS_TOKEN not defined")
		os.Exit(2)
	}

	filterBranch := c.String("branch")
	path := c.String("path")

	client := azuredevops.NewClient(azureDevOpsAccount, azureDevOpsProject, azureDevOpsToken)
	client.UserAgent = "frost/go-azuredevops"

	buildDefOpts := azuredevops.BuildDefinitionsListOptions{Path: "\\" + path}
	definitions, err := client.BuildDefinitions.List(&buildDefOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get a list of build definitions: %v", err)
		os.Exit(2)
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

			for _, branchName := range strings.Split(filterBranch, ",") {
				builds, err := getBuildsForBranch(client, definition.ID, branchName)
				if err != nil {
					fmt.Printf("unable to get builds for definition %s: %v", definition.Name, err)
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

	renderAzureDevOpsBuilds(builds, len(builds), ".*")
}

func getBuildsForBranch(client *azuredevops.Client, defID int, branchName string) ([]azuredevops.Build, error) {
	buildOpts := azuredevops.BuildsListOptions{Definitions: strconv.Itoa(defID), Branch: "refs/heads/" + branchName, Count: 1}
	build, err := client.Builds.List(&buildOpts)
	return build, err
}

func renderAzureDevOpsBuilds(builds []azuredevops.Build, count int, filterBranch string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "", "Name", "Branch", "Build", "Finished")
	for index := 0; index < count; index++ {
		build := builds[index]
		name := build.Definition.Name
		result := build.Result
		status := build.Status
		buildNo := build.BuildNumber
		branch := build.Branch

		// Deal with date formatting for the finish time
		finish, err := time.Parse(time.RFC3339, builds[index].FinishTime)
		finishAt := finish.Format(appDateTimeFormat)
		if err != nil {
			finishAt = builds[index].FinishTime
		}

		// Filter on branches
		matched, _ := regexp.MatchString(".*"+filterBranch+".*", branch)
		if matched == false {
			continue
		}

		if status == "inProgress" {
			result = appProgress
		} else if status == "notStarted" {
			result = appPending
		} else {
			if result == "failed" {
				result = appFailure
			} else {
				result = appSuccess
			}
		}

		fmt.Fprintf(w, "%s \t%s\t%s\t%s\t%s\n", result, name, branch, buildNo, finishAt)
	}

	w.Flush()
}
