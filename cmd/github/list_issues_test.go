package github_test

import (
	"testing"

	"github.com/benmatselby/frost/cmd/github"

	"github.com/benmatselby/frost/test"
	"github.com/spf13/cobra"
)

func TestNewListIssuesCommand(t *testing.T) {
	ctx, client := github.NewClient()

	cmd := github.NewListIssuesCommand(ctx, client)

	expected := &cobra.Command{
		Use:   "list-issues",
		Short: "List all the Issues",
	}

	test.Command(t, cmd, expected)
}
