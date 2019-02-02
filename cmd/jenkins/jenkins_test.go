package jenkins_test

import (
	"testing"

	"github.com/benmatselby/frost/test"

	"github.com/benmatselby/frost/cmd/jenkins"
	"github.com/spf13/cobra"
)

func TestNewJenkinsCommand(t *testing.T) {
	cmd := jenkins.NewJenkinsCommand()

	expected := &cobra.Command{
		Use:   "jenkins",
		Short: "Jenkins related commands",
	}

	test.Command(t, cmd, expected)
}
