# AbyleBotter

[![Build Status](https://travis-ci.org/torlenor/abylebotter.svg?branch=master)](https://travis-ci.org/torlenor/abylebotter)
[![Coverage Status](https://coveralls.io/repos/github/torlenor/AbyleBotter/badge.svg?branch=master)](https://coveralls.io/github/torlenor/AbyleBotter?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/torlenor/AbyleBotter)](https://goreportcard.com/report/github.com/torlenor/AbyleBotter)
[![Docker](https://img.shields.io/docker/pulls/hpsch/abylebotter.svg)](https://hub.docker.com/r/hpsch/abylebotter/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

## Description

This is AbyleBotter, an extensible Chat Bot for various platforms.

At the moment the Bot is in a proof of concept/API/interface development phase with very limited functional use.

## Supported platforms

These platforms are current supported (at least with the functionality to send and receive messages):

### Discord

### Matrix

### Mattermost

### Slack

## Plugins

Plugins are used to implement actual functionality of AbyleBotter. They serve as handlers of received messages and can send messages over the Bot to the platform.
In the future it is planed to support external Plugins via gRPC.

Currently these Plugins are part of AbyleBotter:

### EchoPlugin

The EchoPlugin echos back all messages it received to the sender of the message. It listens to messages which start with `!echo ` followed by text.

#### Configuration Options

- **onlywhispers**: When set to true the EchoPlugin only echos in whispers (when supported by the used Bot)

### HTTPPingPlugin

The HTTPPingPlugin listens to messages starting with `!httpping ` followed by an URL. If the URL is valid it will try to contact the server and reports back the request duration or an error if the URL was not reachable.

### RandomPlugin

This is a classic "Roll/Random" plugin which sends back a random number in the range [0,100] when it receives the `!roll` command. When it receives `!roll {PositiveNumber}` instead, it returns a random number in the range [0, {PositiveNumber}].

### SendMessagePlugin

A small example which shows how one could implement a JSON interface to allow sending messages via an HTTP interface to AbyleBotter.

## Releases

For releases binaries for Linux, Windows and Mac are provided. Check out the respective section on GitHub.

## How to build from source

### Requirements

- Go >= 1.12

### Building

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
- Mattermost: For Mattermost a username and password with the necessary rights on the specified server is enough.
- Slack: The bot as to be added to the workspace and a token has to be generated.

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

where _/path/to/config/file.toml_ has to be replaced with the path to your config file (see above for remarks about bot accounts for the various services).
