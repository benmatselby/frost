package build_test

import (
	"testing"

	"github.com/benmatselby/frost/cmd/ado/build"
	"github.com/benmatselby/frost/test"
	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestNewAdoBuildListCommand(t *testing.T) {
	client := azuredevops.NewClient(
		viper.GetString("AZURE_DEVOPS_ACCOUNT"),
		viper.GetString("AZURE_DEVOPS_PROJECT"),
		viper.GetString("AZURE_DEVOPS_TOKEN"),
	)

	cmd := build.NewAdoBuildListCommand(client)

	expected := &cobra.Command{
		Use:   "list",
		Short: "List all the builds",
	}

	test.Command(t, cmd, expected)
}
