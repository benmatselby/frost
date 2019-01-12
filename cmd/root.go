package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/benmatselby/frost/cmd/ado"
	"github.com/benmatselby/frost/cmd/jenkins"
	"github.com/benmatselby/frost/cmd/travis"
	"github.com/benmatselby/frost/version"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	initConfig()

	cmd := NewRootCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewRootCommand initialises any configuration required
func NewRootCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "frost",
		Short:   "CLI application for retrieving data from the 🌍",
		Version: version.GITCOMMIT,
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.frost/config.yaml)")

	cmd.AddCommand(
		NewVersionCommand(),
		ado.NewAdoCommand(),
		jenkins.NewJenkinsCommand(),
		travis.NewTravisCommand(),
	)
	return cmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".frost" (without extension).
		path := strings.Join([]string{home, ".frost"}, "/")
		viper.AddConfigPath(path)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
