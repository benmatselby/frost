package main

import (
	"fmt"
	"os"

	"github.com/benmatselby/frost/version"
	"github.com/urfave/cli"
)

var (
	vstsAccount string
	vstsProject string
	vstsToken   string

	travisToken string
	travisOwner string
)

const (
	appDateFormat     string = "02-01-2006"
	appDateTimeFormat string = "02-01-2006 15:04"
	appSuccess        string = "✅"
	appFailure        string = "❌"
	appPending        string = "🗂"
	appProgress       string = "🏗"
	appStale          string = "🕳"
	appUnknown        string = "❓"
)

func loadEnvironmentVars() error {
	vstsAccount = os.Getenv("VSTS_ACCOUNT")
	vstsProject = os.Getenv("VSTS_PROJECT")
	vstsToken = os.Getenv("VSTS_TOKEN")

	travisToken = os.Getenv("TRAVIS_CI_TOKEN")
	travisOwner = os.Getenv("TRAVIS_CI_OWNER")

	if vstsAccount == "" || vstsProject == "" || vstsToken == "" {
		return fmt.Errorf("The environment variables are not all set")
	}

	return nil
}

func usage(withError bool) string {
	usage := `
,---.,---.    .---.    .---.  _______
| .-'| .-.\  / .-. )  ( .-._)|__   __|
|  -.|  -'/  | | |(_)(_) \     )| |
| .-'|   (   | | | | _  \ \   (_) |
| |  | |\ \  \  -' /( -'  )    | |
)\|  |_| \)\  )---'   ----'     -'
(__)      (__)(_)

Inspector Jack Frost gets build data out of various build systems into the terminal, where we belong...
`

	if withError {
		usage = usage + `

In order for inspector jack frost to investigate, you need to define the following environment variables:

* VSTS_ACCOUNT = %s
* VSTS_PROJECT = %s
* VSTS_TOKEN   = %s
`
	}
	return usage
}

func main() {
	err := loadEnvironmentVars()
	if err != nil {
		fmt.Fprintln(os.Stderr, usage(true))
		os.Exit(2)
	}

	app := cli.NewApp()
	app.Name = "frost"
	app.Author = "@benmatselby"
	app.Usage = usage(false)
	app.Version = version.GITCOMMIT
	app.Commands = []cli.Command{
		{
			Name:    "travis",
			Aliases: []string{"t"},
			Usage:   "Build data from TravisCI",
			Subcommands: []cli.Command{
				{
					Name:    "overview",
					Action:  travisListBuildOverview,
					Aliases: []string{"o"},
					Usage:   "Get the overview of builds",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "owner", Value: os.Getenv("TRAVIS_CI_OWNER"), Usage: "The owner"},
						cli.StringFlag{Name: "branch", Value: "master", Usage: "Filter by branch name"},
					},
				},
			},
		}, {
			Name:    "vsts",
			Aliases: []string{"v"},
			Usage:   "Build data from the VSTS system",
			Subcommands: []cli.Command{
				{
					Name:    "overview",
					Action:  vstsListBuildOverview,
					Aliases: []string{"o"},
					Usage:   "Get the overview of builds",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "path", Value: os.Getenv("VSTS_TEAM"), Usage: "Build definition path"},
						cli.StringFlag{Name: "branch", Value: "master", Usage: "Filter by branch name"},
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
