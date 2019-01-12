package ado

import (
	"github.com/benmatselby/frost/cmd/ado/build"
	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewAdoCommand will add an `ado` command which will add sub commands
func NewAdoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ado",
		Short: "Azure DevOps related commands",
	}

	client := azuredevops.NewClient(
		viper.GetString("AZURE_DEVOPS_ACCOUNT"),
		viper.GetString("AZURE_DEVOPS_PROJECT"),
		viper.GetString("AZURE_DEVOPS_TOKEN"),
	)
	client.UserAgent = "frost/go-azuredevops"

	cmd.AddCommand(
		build.NewAdoBuildListCommand(client),
	)

	return cmd
}
