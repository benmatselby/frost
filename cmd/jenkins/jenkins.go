package jenkins

import (
	"github.com/benmatselby/frost/cmd/jenkins/build"
	"github.com/benmatselby/frost/cmd/jenkins/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewJenkinsCommand will add an `jenkins` command which will add sub commands
func NewJenkinsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jenkins",
		Short: "Jenkins related commands",
	}

	client := client.New(
		viper.GetString("JENKINS_URL"),
		viper.GetString("JENKINS_USERNAME"),
		viper.GetString("JENKINS_PASSWORD"),
	)

	cmd.AddCommand(
		build.NewJenkinsBuildListCommand(client),
	)

	return cmd
}
