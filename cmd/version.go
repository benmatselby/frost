package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/benmatselby/frost/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand adds a `version` command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version information",

		Run: func(cmd *cobra.Command, args []string) {
			DisplayVersion(os.Stdout)
		},
	}

	return cmd
}

// DisplayVersion will display the current version of the application
func DisplayVersion(w io.Writer) {
	v := version.GITCOMMIT
	fmt.Fprintf(w, "version: %s\n", v)
}
