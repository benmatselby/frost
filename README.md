# Frost

[![Build Status](https://travis-ci.org/benmatselby/frost.png?branch=master)](https://travis-ci.org/benmatselby/frost)
[![codecov](https://codecov.io/gh/benmatselby/frost/branch/master/graph/badge.svg)](https://codecov.io/gh/benmatselby/frost)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/frost?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/frost)

_Inspector Jack Frost_

CLI application for getting certain kinds of data out of various build and work management tools. It currently supports

- [Jenkins](http://jenkins.io) - Build information
- [TravisCI](https://travis-ci.org) - Build information
- [Azure DevOps](https://azure.microsoft.com/en-us/solutions/devops/) - Build information

```output
CLI application for retrieving data from the 🌍

Usage:
  frost [command]

Available Commands:
  ado         Azure DevOps related commands
  help        Help about any command
  jenkins     Jenkins related commands
  travis      TravisCI related commands
  version     Show the version information

Flags:
      --config string   config file (default is $HOME/.frost/config.yaml)
  -h, --help            help for frost

Use "frost [command] --help" for more information about a command.
```

## Requirements

If you are wanting to build and develop this, you will need the following items installed. If, however, you just want to run the application I recommend using the docker container (See below)

- Go version 1.11+

## Configuration

You will need the following environment variables defining depending on which services you want to use:

```shell
export AZURE_DEVOPS_ACCOUNT=""
export AZURE_DEVOPS_PROJECT=""
export AZURE_DEVOPS_TOKEN=""

export TRAVIS_CI_OWNER=""
export TRAVIS_CI_TOKEN=""

export JENKINS_URL=""
export JENKINS_USERNAME=""
export JENKINS_PASSWORD=""
# This is only required, if you want to get an overview of Jenkins from a defined "view".
# If this is not specified, it gets all jobs
export JENKINS_VIEW=""
```

You can also define ~/.frost/config.yml which has various settings.

## Installation via Git

```shell
git clone git@github.com:benmatselby/frost.git
cd frost
make all
./frost --help
```

## Installation via Docker

Other than requiring [docker](http://docker.com) to be installed, there are no other requirements to run the application this way

```shell
$ docker run \
    --rm \
    -t \
    -eAZURE_DEVOPS_ACCOUNT \
    -eAZURE_DEVOPS_PROJECT \
    -eAZURE_DEVOPS_TOKEN \
    -eTRAVIS_CI_OWNER \
    -eTRAVIS_CI_TOKEN \
    -eJENKINS_URL \
    -eJENKINS_USERNAME \
    -eJENKINS_PASSWORD \
    benmatselby/frost "$@"
```

## Bash completion

If you would like bash completion support, then run:

```shell
frost completion
```

This will generate a bash completion script in `/tmp/frost.sh`. You simply need to move this into your `bash_completion.d` folder. On the Mac, this is likely to be `/usr/local/etc/bash_completion.d/`.
