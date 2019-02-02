package build_test

import (
	"testing"

	"github.com/benmatselby/frost/cmd/jenkins/build"
	"github.com/benmatselby/frost/cmd/jenkins/client"
	"github.com/benmatselby/frost/test"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestNewJenkinsBuildListCommand(t *testing.T) {
	client := client.New(
		viper.GetString("JENKINS_URL"),
		viper.GetString("JENKINS_USERNAME"),
		viper.GetString("JENKINS_PASSWORD"),
	)
	cmd := build.NewJenkinsBuildListCommand(client)

	expected := &cobra.Command{
		Use:   "list",
		Short: "List all the builds",
	}

	test.Command(t, cmd, expected)
}
