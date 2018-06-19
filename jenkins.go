package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/jessfraz/tdash/jenkins"
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

	// Initialize the jenkins api client
	jenkinsClient := jenkins.New(jenkinsURL, jenkinsUsername, jenkinsPassword)

	// Get all the jobs
	jobs, err := jenkinsClient.GetJobs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getting jenkins jobs failed: %v", err)
		os.Exit(2)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\t%s\n", "", "Name", "Finished")

	for _, job := range jobs {
		if job.LastBuild.Result == "" {
			// Then the job is currently running.
			job.LastBuild.Result = "RUNNING"
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

		fmt.Fprintf(w, "%s \t%s\t%s\n", result, job.DisplayName, finishedAt)
	}

	w.Flush()
}
