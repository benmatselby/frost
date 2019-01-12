package build

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/frost/cmd/jenkins/client"
	"github.com/benmatselby/frost/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewJenkinsBuildListCommand will add a `jenkins build list` command which
// is responsible for showing a list of build
func NewJenkinsBuildListCommand(client *client.Client) *cobra.Command {
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
func ListBuilds(client *client.Client, w io.Writer) error {
	jobs, err := client.GetJobs(viper.GetString("JENKINS_VIEW"))
	if err != nil {
		return err
	}

	tr := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(tr, "%s\t%s\t%s\t%s\n", "", "Name", "No.", "Finished")

	for _, job := range jobs {
		if job.LastBuild.Result == "" && job.LastBuild.Timestamp != 0 {
			job.LastBuild.Result = "RUNNING"
		}

		if job.LastBuild.Result == "" && job.LastBuild.Timestamp == 0 {
			job.LastBuild.Result = "WAITING"
		}

		finishedAt := time.Unix(0, int64(time.Millisecond)*job.LastBuild.Timestamp).Format(ui.AppDateTimeFormat)

		result := ""
		if job.LastBuild.Result == "FAILURE" {
			result = ui.AppFailure
		} else if job.LastBuild.Result == "SUCCESS" {
			result = ui.AppSuccess
		} else if job.LastBuild.Result == "RUNNING" {
			result = ui.AppProgress
		} else {
			result = ui.AppStale
		}

		fmt.Fprintf(tr, "%s \t%s\t%v\t%s\n", result, job.DisplayName, job.LastBuild.Number, finishedAt)
	}

	tr.Flush()
	return nil
}
