package ado_test

import (
	"testing"

	"github.com/benmatselby/frost/cmd/ado"
	"github.com/benmatselby/frost/test"
	"github.com/spf13/cobra"
)

func TestNewAdoCommand(t *testing.T) {
	cmd := ado.NewAdoCommand()

	expected := &cobra.Command{
		Use:   "ado",
		Short: "Azure DevOps related commands",
	}

	test.Command(t, cmd, expected)
}
