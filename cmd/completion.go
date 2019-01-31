package cmd

import (
	"github.com/spf13/cobra"
)

// NewCompletionCommand creates a new `completion` command
func NewCompletionCommand(root *cobra.Command) *cobra.Command {
	location := "/tmp/frost.sh"
	cmd := &cobra.Command{
		Use:    "completion",
		Short:  "Generates the bash completion script: " + location,
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			root.GenBashCompletionFile(location)
		},
	}
	return cmd
}
