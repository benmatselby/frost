# Frost

[![Build Status](https://travis-ci.org/benmatselby/frost.png?branch=master)](https://travis-ci.org/benmatselby/frost)
[![codecov](https://codecov.io/gh/benmatselby/frost/branch/master/graph/badge.svg)](https://codecov.io/gh/benmatselby/frost)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/frost?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/frost)

_Inspector Jack Frost_

CLI application for getting build information out of various build systems. It currently supports

* [Jenkins](http://jenkins.io)
* [TravisCI](https://travis-ci.org)
* Visual Studio Team Services

## Requirements

If you are wanting to build and develop this, you will need the following items installed. If, however, you just want to run the application I recommend using the docker container (See below)

* Go version 1.10+
* [Dep installed](https://github.com/golang/dep)

## Configuration

You will need the following environment variables defining depending on which services you want to use:

```bash
$ export VSTS_ACCOUNT=""
$ export VSTS_PROJECT=""
$ export VSTS_TOKEN=""

$ export TRAVIS_CI_OWNER=""
$ export TRAVIS_CI_TOKEN=""

$ export JENKINS_URL=""
$ export JENKINS_USERNAME=""
$ export JENKINS_PASSWORD=""
$ export JENKINS_VIEW="" # This is only required, if you want to get an overview of Jenkins from a defined "view"
```

## Installation via Git

```bash
$ git clone git@github.com:benmatselby/frost.git
$ cd frost
$ make all
$ ./frost --help
```

## Installation via Docker

Other than requiring [docker](http://docker.com) to be installed, there are no other requirements to run the application this way

```bash
$ docker run \
    --rm \
    -t \
    -eVSTS_ACCOUNT \
    -eVSTS_PROJECT \
    -eVSTS_TOKEN \
    -eTRAVIS_CI_OWNER \
    -eTRAVIS_CI_TOKEN \
    -eJENKINS_URL \
    -eJENKINS_USERNAME \
    -eJENKINS_PASSWORD \
    benmatselby/frost "$@"
```
