package test

import (
	"testing"

	"github.com/spf13/cobra"
)

// Command knows how to assert a cobra.Command
func Command(t *testing.T, command, expected *cobra.Command) {
	if command.Use != expected.Use {
		t.Fatalf("expected use %s; got %s", expected.Use, command.Use)
	}

	if command.Short != expected.Short {
		t.Fatalf("expected use %s; got %s", expected.Short, command.Short)
	}
}
