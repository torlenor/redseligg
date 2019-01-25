# AbyleBotter

[![Build Status](https://travis-ci.org/torlenor/abylebotter.svg?branch=master)](https://travis-ci.org/torlenor/abylebotter)
[![Coverage Status](https://coveralls.io/repos/github/torlenor/AbyleBotter/badge.svg?branch=master)](https://coveralls.io/github/torlenor/AbyleBotter?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/torlenor/AbyleBotter)](https://goreportcard.com/report/github.com/torlenor/AbyleBotter)
[![Docker](https://img.shields.io/docker/pulls/hpsch/abylebotter.svg)](https://hub.docker.com/r/hpsch/abylebotter/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

## Description

This is AbyleBotter, an extensible Chat Bot for Discord and Matrix (+ in the future other chat platforms).

At the moment the Bot is in a proof of concept/API/interface development phase with very limited functional use.

## Releases

For releases binaries for Linux, Windows and Mac are provided. Check out the respective section on GitHub.

## How to build from source

Clone the sources from GitHub and compile AbyleBotter with

```
make deps
make
```

and optionally

```
make test
```

to run tests.

## How to run it

Independent of the way you obtain it, you have to configure the bot first and it is necessary to have a registered bot account for the service you want to use. 

- Discord: Please take a look at https://discordapp.com/developers/docs/intro on how to set up a bot user and generate the required authentication token.
- Matrix: For Matrix it is simpler, just create a user for the bot on your preferred Matrix server.

Then please take a look at the provided example configuration in _config/config.toml_ and adapt it to match your settings.

To start AbyleBotter using the self-built or downloaded binary enter

```
./path/to/abylebotter -c /path/to/config/file.toml
```

The Bot should now connect automatically to the service and should be ready to use.

## Using Docker

Probably the easiest way to try out AbyleBotter is using Docker. To pull the latest version from DockerHub and start it just type

```
docker run --name abylebotter -v /path/to/config/file.toml:/app/config/config.toml:ro hpsch/abylebotter:latest
```

where _/path/to/config/file.toml_ has to be replaced with the path to your config file.
