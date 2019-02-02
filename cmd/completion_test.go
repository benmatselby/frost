package cmd_test

import (
	"testing"

	"github.com/benmatselby/frost/cmd"
	"github.com/benmatselby/frost/test"
	"github.com/spf13/cobra"
)

func TestNewCompletionCommand(t *testing.T) {
	root := &cobra.Command{
		Use: "mock",
	}

	cmd := cmd.NewCompletionCommand(root)

	expected := &cobra.Command{
		Use:   "completion",
		Short: "Generates the bash completion script: /tmp/frost.sh",
	}

	test.Command(t, cmd, expected)
}
