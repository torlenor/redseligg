# AbyleBotter

[![Build Status](https://travis-ci.org/torlenor/AbyleBotter.svg?branch=master)](https://travis-ci.org/torlenor/AbyleBotter)
[![Go Report Card](https://goreportcard.com/badge/github.com/torlenor/AbyleBotter)](https://goreportcard.com/report/github.com/torlenor/AbyleBotter)
[![Docker](https://img.shields.io/docker/pulls/hpsch/abylebotter.svg)](https://hub.docker.com/r/hpsch/abylebotter/)
[![License](https://img.shields.io/badge/license-GPL-blue.svg)](/LICENSE)

## Description

This is AbyleBotter, an extensible Chat Bot for Discord (and in the future other chat platforms).

At the moment the Bot is in a proof of concept/API/interface development phase with very limited functional use.

## How to download/install

Checkout and compile AbyleBotter with

```
$ git clone https://github.com/torlenor/AbyleBotter.git
$ cd AbyleBotter
$ make deps
$ make
```

## How to start

Currently it is necessary to have a registered Bot Account which has already joined some servers. Please take a look at https://discordapp.com/developers/docs/intro on how to set this up. In the end you should have a bot token to use with AbyleBotter.

To start AbyleBotter use when using Linux (bash or similar shell)

```
$ DISCORD_BOT_TOKEN=token ./bin/AbyleBotter
```

The Bot should now connect to the Discord Gateway and should be ready to use. 

## Using Docker

If you want try out AbyleBotter in a Docker container you can pull the latest version with

```
$ docker pull hpsch/abylebotter:latest
```

After the successful download it can be started with

```
$ docker run -i -t --name abylebotter --rm -e DISCORD_BOT_TOKEN=token hpsch/abylebotter
```

where token has to be replaced with your Discord Bot token.
