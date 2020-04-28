# AbyleBotter

[![Build status](https://git.abyle.org/hps/abylebotter/badges/master/pipeline.svg)](https://git.abyle.org/hps/abylebotter/commits/master)
[![Coverage Status](https://git.abyle.org/hps/abylebotter/badges/master/coverage.svg)](https://git.abyle.org/hps/abylebotter/commits/master)
[![Go Report Card](https://goreportcard.com/badge/git.abyle.org/hps/abylebotter)](https://goreportcard.com/report/git.abyle.org/hps/abylebotter)
[![Docker](https://img.shields.io/docker/pulls/hpsch/abylebotter.svg)](https://hub.docker.com/r/hpsch/abylebotter/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)


## Description

This is AbyleBotter, an extensible Chat Bot for various platforms. It is based on a server architecture which can be controlled via a REST API.

## Supported platforms

These platforms are current supported (at least with the functionality to send and receive messages):

### Discord

### Matrix

### Mattermost

### Slack

## Releases

For releases binaries for Linux, Windows and Mac are provided. Check out the respective section on GitHub.

## How to build from source

### Requirements

- Go >= 1.12

### Building

Clone the sources from the repository and compile it with

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

- Discord: Please take a look at https://discordapp.com/developers/docs/intro on how to set up a bot user and generate the required authentication token. Then use the bot OAuth2 authorization link, which can be generated on your applications page at OAuth2 when you select as scope "Bot". Note: This authentication flow is much easier than the normal OAuth2 user challenge and does not require a callback link. For details on that visit https://discordapp.com/developers/docs/topics/oauth2#bot-authorization-flow.
- Matrix: For Matrix it is simpler, just create a user for the bot on your preferred Matrix server.
- Mattermost: For Mattermost a username and password with the necessary rights on the specified server is enough.
- Slack: The bot as to be added to the workspace and a token has to be generated.

The bot configuration can either be stored in a toml file or in a MongoDB. An example for a toml file is provided in this repository in *cfg/bots.toml*.

### Self-built or downloaded release

To start the AbyleBotter BotterInstance using the self-built or downloaded binary enter for use with a TOML config

```bash
BOTTER_BOT_CFG_SOURCE="TOML" BOTTER_BOT_CFG_TOML_FILE=/path/to/config/file.toml ./botterinstance
```

or 

```bash
BOTTER_BOT_CFG_SOURCE="MONGO" BOTTER_BOT_CFG_MONGO_URL="mongodb://user:password@localhost/database" BOTTER_BOT_CFG_MONGO_DB="database" ./botterinstance
```

for use with a MongoDB for the Bot configuration.

### Using Docker

Probably an easiest way to try out AbyleBotter is using Docker. To pull the latest version from DockerHub and start it type

```bash
docker run -d --name abylebotter --env BOTTER_BOT_CFG_SOURCE=TOML --env BOTTER_BOT_CFG_TOML_FILE=/bots.toml -v /path/to/config/file.toml:/bots.toml:ro hpsch/abylebotter:latest
```

or for MongoDB type

```bash
docker run -d --name abylebotter BOTTER_BOT_CFG_SOURCE="MONGO" BOTTER_BOT_CFG_MONGO_URL="mongodb://user:password@localhost/database" BOTTER_BOT_CFG_MONGO_DB="database" hpsch/abylebotter:latest
```

## How to control it

We are providing a command line tool to control a BotterInstance called BotterControl.

### Get all running bots of a BotterInstance

```bash
./bottercontrol -u URL_OF_BOTTER_INSTANCE -c GetBots
```

or

```bash
docker run --net host hpsch/abylebotter:latest /usr/bin/bottercontrol -u URL_OF_BOTTER_INSTANCE -c GetBots
```

### Start a bot on a BotterInstance

```bash
./bottercontrol -u URL_OF_BOTTER_INSTANCE -c StartBot -a BOTID
```

or

```bash
docker run --net host hpsch/abylebotter:latest /usr/bin/bottercontrol -u URL_OF_BOTTER_INSTANCE -c StartBot -a BOTID
```

### Stop a bot on a BotterInstance

```bash
./bottercontrol -u URL_OF_BOTTER_INSTANCE -c StopBot -a BOTID
```

or

```bash
docker run --net host hpsch/abylebotter:latest /usr/bin/bottercontrol -u URL_OF_BOTTER_INSTANCE -c StopBot -a BOTID
```

## Standalone version

We also provide a standalone version which does not depend on a control instance to launch bots, but just starts all enabled bots from the configuration.

You launch the standalone version type

```bash
BOTTER_BOT_CFG_SOURCE="TOML" BOTTER_BOT_CFG_TOML_FILE=/path/to/config/file.toml ./botter
```

or 

```bash
BOTTER_BOT_CFG_SOURCE="MONGO" BOTTER_BOT_CFG_MONGO_URL="mongodb://user:password@localhost/database" BOTTER_BOT_CFG_MONGO_DB="database" ./botter
```

for use with a MongoDB for the Bot configuration.

### Using Docker

To launch the standalone version using Docker type

```bash
docker run -d --name abylebotter --env BOTTER_BOT_CFG_SOURCE=TOML --env BOTTER_BOT_CFG_TOML_FILE=/bots.toml -v /path/to/config/file.toml:/bots.toml:ro hpsch/abylebotter:latest /usr/bin/botter
```

or for MongoDB type

```bash
docker run -d --name abylebotter BOTTER_BOT_CFG_SOURCE="MONGO" BOTTER_BOT_CFG_MONGO_URL="mongodb://user:password@localhost/database" BOTTER_BOT_CFG_MONGO_DB="database" hpsch/abylebotter:latest /usr/bin/botter
```

## Plugins

Plugins are used to implement actual functionality of AbyleBotter. They serve as handlers of received messages and can send messages over the Bot to the platform.
In the future it is planed to support external Plugins via gRPC.

Currently these Plugins are part of AbyleBotter:

### EchoPlugin

The EchoPlugin echos back all messages it received to the sender of the message. It listens to messages which start with `!echo ` followed by text.

#### Configuration Options

- **onlywhispers**: When set to true the EchoPlugin only echos in whispers (when supported by the used Bot)

### GiveawayPlugin

Lets you hold giveaways in your channel and let the bot pick a winner.

#### Configuration options

Example:
```toml
[bots.some_bot.plugins.1]
    type = "giveaway"
    [bots.some_bot.plugins.1.config]
        mods = ["user"]
        onlymods = true
```

When `onlymods` is set to `true`, only the users which are listed in `mods` are allowed to start/end giveaways. Per default everybody is allowed.

#### Starting a giveaway

To start a giveaway in the current channel type

```
!gstart <time> <secretword> [winners] [prize]
```

* `<time>` is the time the giveaway should run. It should include s/m/h to indicate seconds/minutes/hours.
* `<secretword>` is the word the bot should react for the people to participate.
* `[winners]` is the number of winners to pick in the end.
* `[prize]` is the prize the people can win.

Example:

```
!gstart 1m hello 2 Bananas
```

#### Ending a giveaway

To stop a currently running giveaway type `!gend`.

#### Rerolling the winner

Type `!greroll` to pick a new winner from the last ended giveaway.

### HTTPPingPlugin

The HTTPPingPlugin listens to messages starting with `!httpping ` followed by an URL. If the URL is valid it will try to contact the server and reports back the request duration or an error if the URL was not reachable.

### RandomPlugin

This is a classic "Roll/Random" plugin which sends back a random number in the range [0,100] when it receives the `!roll` command. When it receives `!roll {PositiveNumber}` instead, it returns a random number in the range [0, {PositiveNumber}].

### VersionPlugin

The VersionPlugin answers to `!version` with the version of the bot.

### Vote Plugin

Initiate a vote in the channel about arbitrary topics.

#### Configuration options

Example:
```toml
[bots.some_bot.plugins.1]
    type = "vote"
    [bots.some_bot.plugins.1.config]
        mods = ["user"]
        onlymods = true
```

When `onlymods` is set to `true`, only the users which are listed in `mods` are allowed to start/end vots. Per default everybody is allowed.

#### Starting a vote
Type `!vote message` to start the vote. The vote is limited to the channel where you initiate the vote. Per default the options are Yes/No. They can be changed by providing custom options (see below).

#### How to participate in the vote
React with the emoji assigned to the options you want to vote for.

#### Custom voting options
Provide the custom options in square brackets after the message, e.g., 
```
!vote What is the best color? [Red, Green, Blue]
```

#### Ending a vote

Type '!voteend message' to end a vote. No additional choices will be counted. For example to end the vote started above type
```
!voteend What is the best color?
```

#### Deleting a vote

Just delete the vote message.

## Quotes Plugin

Lets users/viewers or mods add quotes and randomly fetch one.

### Configuration options

Example:
```toml
[bots.some_bot.plugins.1]
    type = "quotes"
    [bots.some_bot.plugins.1.config]
        mods = ["user"]
        onlymods = true
```

When `onlymods` is set to `true`, only the users which are listed in `mods` are allowed to perform certain actions. Per default everybody is allowed.

### Adding a quote

To add a quote type

```
!quoteadd <your quote>
```

Example:
```
!quoteadd This is awesome!
```

### Getting a quote

```
!quote
```
will return a random quote.

```
!quote 2
```
will return the 2nd quote in the list.

The output will be similar to

```
123. "This is awesome!" - 2020-4-22, added by SomeBody
```

### Removing a quote

Use 

```
!quoteremove ID
```

, e.g., !quoteremove 123, to remove a quote.

**Note:** When `onlymods` is set to `true` in configuration, only mods are allowed to list all quotes.
