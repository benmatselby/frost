# Frost

[![Build Status](https://travis-ci.org/benmatselby/frost.png?branch=master)](https://travis-ci.org/benmatselby/frost)
[![codecov](https://codecov.io/gh/benmatselby/frost/branch/master/graph/badge.svg)](https://codecov.io/gh/benmatselby/frost)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/frost?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/frost)

_Inspector Jack Frost_

CLI application for getting certain kinds of data out of various build and work management tools. It currently supports

* [Jenkins](http://jenkins.io) - Build information
* [TravisCI](https://travis-ci.org) - Build information
* [Azure DevOps](https://azure.microsoft.com/en-us/solutions/devops/) - Build information
* [GitHub](https://github.com)- Pull Request information

## Requirements

If you are wanting to build and develop this, you will need the following items installed. If, however, you just want to run the application I recommend using the docker container (See below)

* Go version 1.11+

## Configuration

You will need the following environment variables defining depending on which services you want to use:

```bash
export AZURE_DEVOPS_ACCOUNT=""
export AZURE_DEVOPS_PROJECT=""
export AZURE_DEVOPS_TOKEN=""

export TRAVIS_CI_OWNER=""
export TRAVIS_CI_TOKEN=""

export GITHUB_ORG=""
export GITHUB_OWNER=""
export GITHUB_TOKEN=""

export JENKINS_URL=""
export JENKINS_USERNAME=""
export JENKINS_PASSWORD=""
# This is only required, if you want to get an overview of Jenkins from a defined "view".
# If this is not specified, it gets all jobs
export JENKINS_VIEW=""
```

You can also define ~/.frost/config.yml which has various settings.

### Limiting the repos to show Pull Requests for

```yml
github:
  pull_request_repos:
  - my-org/my-repo
  - benmatselby/*
```

## Installation via Git

```bash
git clone git@github.com:benmatselby/frost.git
cd frost
make all
./frost --help
```

## Installation via Docker

Other than requiring [docker](http://docker.com) to be installed, there are no other requirements to run the application this way

```bash
$ docker run \
    --rm \
    -t \
    -eAZURE_DEVOPS_ACCOUNT \
    -eAZURE_DEVOPS_PROJECT \
    -eAZURE_DEVOPS_TOKEN \
    -eGITHUB_ORG \
    -eGITHUB_OWNER \
    -eGITHUB_TOKEN \
    -eTRAVIS_CI_OWNER \
    -eTRAVIS_CI_TOKEN \
    -eJENKINS_URL \
    -eJENKINS_USERNAME \
    -eJENKINS_PASSWORD \
    benmatselby/frost "$@"
```
