package cmd_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/benmatselby/frost/version"

	"github.com/benmatselby/frost/cmd"
	"github.com/benmatselby/frost/test"
	"github.com/spf13/cobra"
)

func TestNewVersionCommand(t *testing.T) {
	cmd := cmd.NewVersionCommand()

	expected := &cobra.Command{
		Use:   "version",
		Short: "Show the version information",
	}

	test.Command(t, cmd, expected)
}

func TestDisplayVersion(t *testing.T) {
	version.GITCOMMIT = "xyz"
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	cmd.DisplayVersion(writer)
	writer.Flush()

	expected := "version: xyz\n"
	if b.String() != expected {
		t.Fatalf("expected '%s'; got '%s'", expected, b.String())
	}
}
