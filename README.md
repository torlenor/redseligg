# Redseligg

[![Build and Test](https://github.com/torlenor/redseligg/workflows/Build%20and%20Test/badge.svg?branch=master)](https://github.com/torlenor/redseligg/actions?query=workflow%3A%22Build+and+Test%22)
[![Coverage Status](https://coveralls.io/repos/github/torlenor/redseligg/badge.svg?branch=master)](https://coveralls.io/github/torlenor/redseligg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/torlenor/redseligg)](https://goreportcard.com/report/github.com/torlenor/redseligg)
[![Docker](https://img.shields.io/docker/pulls/torlenor/redseligg.svg)](https://hub.docker.com/r/torlenor/redseligg/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)


## Description

This is Redseligg, an extensible Chat Bot for various platforms. It is based on a server architecture which can be controlled via a REST API.

## Supported platforms

These platforms are current supported (at least with the functionality to send and receive messages):

### Discord

### Matrix

### Mattermost

### Slack

### Twitch

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
- Twitch: It needs a username for the Twitch account and a list of channels to join. In addition a toen is needed for that user. You can generate one here: https://twitchapps.com/tmi/

The bot configuration can either be stored in a toml file or in a MongoDB. An example for a toml file is provided in this repository in *cfg/bots.toml*.

### Self-built or downloaded release

To start the Redseligg BotterInstance using the self-built or downloaded binary enter for use with a TOML config

```bash
BOTTER_BOT_CFG_SOURCE="TOML" BOTTER_BOT_CFG_TOML_FILE=/path/to/config/file.toml ./botterinstance
```

or 

```bash
BOTTER_BOT_CFG_SOURCE="MONGO" BOTTER_BOT_CFG_MONGO_URL="mongodb://user:password@localhost/database" BOTTER_BOT_CFG_MONGO_DB="database" ./botterinstance
```

for use with a MongoDB for the Bot configuration.

### Using Docker

Probably an easiest way to try out Redseligg is using Docker. To pull the latest version from DockerHub and start it type

```bash
docker run -d --name redseligg --env BOTTER_BOT_CFG_SOURCE=TOML --env BOTTER_BOT_CFG_TOML_FILE=/bots.toml -v /path/to/config/file.toml:/bots.toml:ro hpsch/redseligg:latest
```

or for MongoDB type

```bash
docker run -d --name redseligg BOTTER_BOT_CFG_SOURCE="MONGO" BOTTER_BOT_CFG_MONGO_URL="mongodb://user:password@localhost/database" BOTTER_BOT_CFG_MONGO_DB="database" hpsch/redseligg:latest
```

## How to control it

We are providing a command line tool to control a BotterInstance called BotterControl.

### Get all running bots of a BotterInstance

```bash
./bottercontrol -u URL_OF_BOTTER_INSTANCE -c GetBots
```

or

```bash
docker run --net host hpsch/redseligg:latest /usr/bin/bottercontrol -u URL_OF_BOTTER_INSTANCE -c GetBots
```

### Start a bot on a BotterInstance

```bash
./bottercontrol -u URL_OF_BOTTER_INSTANCE -c StartBot -a BOTID
```

or

```bash
docker run --net host hpsch/redseligg:latest /usr/bin/bottercontrol -u URL_OF_BOTTER_INSTANCE -c StartBot -a BOTID
```

### Stop a bot on a BotterInstance

```bash
./bottercontrol -u URL_OF_BOTTER_INSTANCE -c StopBot -a BOTID
```

or

```bash
docker run --net host hpsch/redseligg:latest /usr/bin/bottercontrol -u URL_OF_BOTTER_INSTANCE -c StopBot -a BOTID
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
docker run -d --name redseligg --env BOTTER_BOT_CFG_SOURCE=TOML --env BOTTER_BOT_CFG_TOML_FILE=/bots.toml -v /path/to/config/file.toml:/bots.toml:ro hpsch/redseligg:latest /usr/bin/botter
```

or for MongoDB type

```bash
docker run -d --name redseligg BOTTER_BOT_CFG_SOURCE="MONGO" BOTTER_BOT_CFG_MONGO_URL="mongodb://user:password@localhost/database" BOTTER_BOT_CFG_MONGO_DB="database" hpsch/redseligg:latest /usr/bin/botter
```

## Storage

Some plugins can use a storage to store permanent data. Currently we are supporting MongoDB and SQLite3 as a storage backend.

### MongoDB

To configure a MongoDB storage, add a section of the form

```yaml
    [bots.slack.storage]
      storage = "mongo"
      [bots.slack.storage.config]
        URL = "mongodb://user1:test@localhost/testdb" # URL to connect to
        Database = "testdb" # Database to use
```

to the bot for which you want to enable the storage (in the example above for a bot called 'slack'). The plugins for this bot will automatically use that storage.

### SQLite3

To configure a SQLite3 storage, add a section of the form

```yaml
    [bots.slack.storage]
      storage = "sqlite"
      [bots.slack.storage.config]
        Database = "/tmp/database.db" # Database file to use
```

to the bot for which you want to enable the storage (in the example above for a bot called 'slack'). The plugins for this bot will automatically use that storage.

## Plugins

Plugins are used to implement actual functionality of Redseligg. They serve as handlers of received messages and can send messages over the Bot to the platform.
In the future it is planed to support external Plugins via gRPC.

Currently these Plugins are part of Redseligg:

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
123. "This is awesome!" - 2020-4-22, added by Somebody
```

### Removing a quote

Use 

```
!quoteremove ID
```

, e.g., `!quoteremove 123`, to remove a quote.

**Note:** When `onlymods` is set to `true` in configuration, only mods are allowed to list all quotes.

## Timed Messages Plugin

Post messages automatically at a given interval.

### Configuration options

Example:
```toml
[bots.some_bot.plugins.1]
    type = "timedmessages"
    [bots.some_bot.plugins.1.config]
        mods = ["user"]
        onlymods = true
```

When `onlymods` is set to `true`, only the users which are listed in `mods` are allowed to add or removed timed messages. Per default everybody is allowed.

### Adding a timed message

To add a timed message type

```
!tm add <interval> <your message>
```

Example:
```
!tm add 1m This is awesome!
```

### Removing a timed message

Use
```
!tm remove <interval> <your message>
```
, e.g., `!tm remove 1m This is awesome!`, to remove one message with a specific interval and text.

or use
```
!tm remove all <your message>
```
, e.g., `!tm remove all This is awesome!`, to remove all message with a specific text, regardless of their interval.

## Custom Commands Plugin

Allow mods to add custom commands which will return text.

### Configuration options

Example:
```toml
[bots.some_bot.plugins.1]
    type = "customcommands"
    [bots.some_bot.plugins.1.config]
        mods = ["user"]
        onlymods = true
```

When `onlymods` is set to `true`, only the users which are listed in `mods` are allowed to add or removed custom commands. Per default everybody is allowed.

### Adding a custom command

To add a custom command type

```
!customcommand add <customCommand> <your message>
```

Example:
```
!customcommand hello Hi there!
```

When a user then types `!hello` in chat the plugin will answer with `Hi there!`.

### Removing a custom command

Use
```
!customcommand remove <customCommand>
```
, e.g., `!tm remove hello`, to remove the custom command.
