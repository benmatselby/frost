package travis_test

import (
	"testing"

	"github.com/benmatselby/frost/cmd/travis"
	"github.com/benmatselby/frost/test"
	"github.com/spf13/cobra"
)

func TestNewTravisCommand(t *testing.T) {
	cmd := travis.NewTravisCommand()

	expected := &cobra.Command{
		Use:   "travis",
		Short: "TravisCI related commands",
	}

	test.Command(t, cmd, expected)
}
