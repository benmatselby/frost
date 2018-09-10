package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/frost/jenkins"
	"github.com/urfave/cli"
)

func jenkinsListBuildOverview(c *cli.Context) {
	if len(jenkinsURL) <= 0 {
		fmt.Fprintf(os.Stderr, "os.Env JENKINS_URL not defined")
		os.Exit(2)
	}

	if len(jenkinsUsername) <= 0 {
		fmt.Fprintf(os.Stderr, "os.Env JENKINS_USERNAME not defined")
		os.Exit(2)
	}

	if len(jenkinsPassword) <= 0 {
		fmt.Fprintf(os.Stderr, "os.Env JENKINS_PASSWORD not defined")
		os.Exit(2)
	}

	client := jenkins.New(jenkinsURL, jenkinsUsername, jenkinsPassword)

	jobs, err := client.GetJobs(jenkinsView)
	if err != nil {
		fmt.Fprintf(os.Stderr, "getting jenkins jobs failed: %v", err)
		os.Exit(2)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "", "Name", "No.", "Finished")

	for _, job := range jobs {
		if job.LastBuild.Result == "" {
			// Assumption made here is that this is a folder/pipline entry
			// with no useful information to render
			continue
		}

		finishedAt := time.Unix(0, int64(time.Millisecond)*job.LastBuild.Timestamp).Format(appDateTimeFormat)

		result := ""
		if job.LastBuild.Result == "FAILURE" {
			result = appFailure
		} else if job.LastBuild.Result == "SUCCESS" {
			result = appSuccess
		} else if job.LastBuild.Result == "RUNNING" {
			result = appProgress
		} else {
			result = appStale
		}

		fmt.Fprintf(w, "%s \t%s\t%v\t%s\n", result, job.DisplayName, job.LastBuild.Number, finishedAt)
	}

	w.Flush()
}
