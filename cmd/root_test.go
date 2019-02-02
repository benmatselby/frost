package cmd_test

import (
	"testing"

	"github.com/benmatselby/frost/test"

	"github.com/benmatselby/frost/cmd"
	"github.com/spf13/cobra"
)

func TestNewRootCommand(t *testing.T) {
	cmd := cmd.NewRootCommand()

	expected := &cobra.Command{
		Use:   "frost",
		Short: "CLI application for retrieving data from the 🌍",
	}

	test.Command(t, cmd, expected)
}
