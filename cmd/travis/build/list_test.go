package build_test

import (
	"testing"

	travis "github.com/Ableton/go-travis"
	"github.com/benmatselby/frost/cmd/travis/build"
	"github.com/benmatselby/frost/test"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestNewAdoBuildListCommand(t *testing.T) {
	client := travis.NewClient(
		travis.TRAVIS_API_DEFAULT_URL,
		viper.GetString("TRAVIS_CI_TOKEN"),
	)

	cmd := build.NewTravisBuildListCommand(client)

	expected := &cobra.Command{
		Use:   "list",
		Short: "List all the builds",
	}

	test.Command(t, cmd, expected)
}
