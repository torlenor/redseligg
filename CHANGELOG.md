# Change Log

## [0.0.6](https://github.com/torlenor/redseligg/releases/tag/0.0.6) (TBD)

**Implemented enhancements:**

**New storage support:**

**New platforms:**

**New plugins:**

- Archive plugin: Stores all messages with their timestamps in the storage.
- RSS Plugin: Subscribe to RSS feeds.

## [0.0.5](https://github.com/torlenor/redseligg/releases/tag/0.0.5) (2020-05-22)

**Implemented enhancements:**

- Added storage system to Plugin API.
- Added requested/provided features check between plugin and platform.
- Added OnCommand() hook from Bot to Plugin which is called when a messages with a previously registered command is received.
- Added !help command which lists all the available commands (including commands from the CustomCommandsPlugin).
- Call-prefix for commands is now customizable via config (default is still '!').

**New storage support:**

- Storage memory added.
- Storage mongo added.
- Storage sqlite added.

**New platforms:**

- Twitch Chat: Redseligg now supports Twitch Chat for sending/receiving messages.

**New plugins:**

- Custom Commands Plugin: Allow mods to add custom commands which will return text.
- Quotes Plugin: Lets users/viewers or mods add quotes and randomly fetch one.
- Timed Messages Plugin: Post messages automatically at a given interval.

## [0.0.4](https://github.com/torlenor/redseligg/releases/tag/0.0.4) (2020-04-21)

**Implemented enhancements:**

- Added GetVersion() to the Bot/Plugin API.
- Added OnRun() hook from Bot to Plugin which is called when the Bot is ready to serve the plugins.
- Added OnStop() hook from Bot to Plugin which is called when the Bot is shutting down.
- Added OnReactionAdded, OnReactionRemoved hooks from Bot to Plugin.
- Add a standalone binary which automatically starts all enabled bots from the provided config file / MongoDB.

**New plugins:**

- Version Plugin: To the command *!version* the plugin will answer with the version of the Bot.
- Giveaway Plugin: Lets you hold giveaways in your channel and let the bot pick a winner.
- Vote Plugin: Initiate a vote in the channel about arbitrary topics.

**Platform specific changes:**

- Discord: Initial mapping of Redseligg to Discord emojis.
- Discord: Initial converter between Redseligg and Discord text format.
- Discord: Support for updating/deleting message.
- Discord: Support for receiving reactions to messages.
- Slack: Support for updating messages.
- Slack: Support for receiving reactions to messages.

**Closed issues:**
- Discord: User.Name shall be User#1234 instead of just User
- Discord: Bot should be able to recover in case of errors
- Discord: Make Callback URL configurable

## [0.0.3](https://github.com/torlenor/redseligg/releases/tag/0.0.3) (2020-03-31)

*Major rework of the whole project to a modern and more Go-like structure.*

**Implemented enhancements:**

- Migrated to a server architecture where Redseligg is controlled via REST API
- Introduce the command line tool BotterControl to control a Redseligg instance
- Cleaner plugin interface and therefore much easier to implement new plugins

**Fixed bugs:**

- Various bugs fixed all over the place

## [0.0.2](https://github.com/torlenor/redseligg/releases/tag/0.0.2) (2018-10-06)

*First release.*
