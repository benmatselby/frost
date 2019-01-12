package travis

import (
	travis "github.com/Ableton/go-travis"
	"github.com/benmatselby/frost/cmd/travis/build"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewTravisCommand will add a `travis` command which will add sub commands
func NewTravisCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "travis",
		Short: "TravisCI related commands",
	}

	client := travis.NewClient(
		travis.TRAVIS_API_DEFAULT_URL,
		viper.GetString("TRAVIS_CI_TOKEN"),
	)

	cmd.AddCommand(
		build.NewTravisBuildListCommand(client),
	)

	return cmd
}
