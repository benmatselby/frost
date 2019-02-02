package github_test

import (
	"testing"

	"github.com/benmatselby/frost/cmd/github"

	"github.com/benmatselby/frost/test"
	"github.com/spf13/cobra"
)

func TestNewPullRequestCommand(t *testing.T) {
	ctx, client := github.NewClient()

	cmd := github.NewPullRequestCommand(ctx, client)

	expected := &cobra.Command{
		Use:   "list-pr",
		Short: "List all the Pull Requests",
	}

	test.Command(t, cmd, expected)
}
