package main

import (
	"os"

	"github.com/benmatselby/frost/version"
	"github.com/urfave/cli"
)

var (
	vstsAccount string
	vstsProject string
	vstsTeam    string
	vstsToken   string

	travisOwner string
	travisToken string

	jenkinsURL      string
	jenkinsUsername string
	jenkinsPassword string
	jenkinsView     string
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

func loadEnvironmentVars() {
	vstsAccount = os.Getenv("VSTS_ACCOUNT")
	vstsProject = os.Getenv("VSTS_PROJECT")
	vstsTeam = os.Getenv("VSTS_TEAM")
	vstsToken = os.Getenv("VSTS_TOKEN")

	travisOwner = os.Getenv("TRAVIS_CI_OWNER")
	travisToken = os.Getenv("TRAVIS_CI_TOKEN")

	jenkinsURL = os.Getenv("JENKINS_URL")
	jenkinsUsername = os.Getenv("JENKINS_USERNAME")
	jenkinsPassword = os.Getenv("JENKINS_PASSWORD")
	jenkinsView = os.Getenv("JENKINS_VIEW")
}

func usage() string {
	usage := `

,---.,---.    .---.    .---.  _______
| .-'| .-.\  / .-. )  ( .-._)|__   __|
|  -.|  -'/  | | |(_)(_) \     )| |
| .-'|   (   | | | | _  \ \   (_) |
| |  | |\ \  \  -' /( -'  )    | |
)\|  |_| \)\  )---'   ----'     -'
(__)      (__)(_)

Inspector Jack Frost gets build data out of various build systems into the terminal, where we belong...

In order for inspector jack frost to investigate, you need to define the following environment variables, depending on
which systems you want to communicate with:

* VSTS_ACCOUNT
* VSTS_PROJECT
* VSTS_TEAM
* VSTS_TOKEN

* TRAVIS_CI_OWNER
* TRAVIS_CI_TOKEN

* JENKINS_URL
* JENKINS_USERNAME
* JENKINS_PASSWORD
* JENKINS_VIEW
`

	return usage
}

func main() {
	loadEnvironmentVars()

	app := cli.NewApp()
	app.Name = "frost"
	app.Author = "@benmatselby"
	app.Usage = usage()
	app.Version = version.GITCOMMIT
	app.Commands = []cli.Command{
		{
			Name:    "jenkins",
			Aliases: []string{"j"},
			Usage:   "Build data from Jenkins",
			Subcommands: []cli.Command{
				{
					Name:    "overview",
					Action:  jenkinsListBuildOverview,
					Aliases: []string{"o"},
					Usage:   "Get the overview of builds",
				},
			},
		},
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
		},
		{
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
