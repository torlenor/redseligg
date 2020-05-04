# Change Log

## [0.0.5](https://git.abyle.org/hps/abylebotter/-/tree/0.0.5) (TBD)

**Implemented enhancements:**

- Added storage system to Plugin API.
- Added requested/provided features check between plugin and platform.

**New storage support:**

- Storage memory added.
- Storage mongo added.

** New platforms:**

- Twitch Chat: AbyleBotter now supports Twitch Chat for sending/receiving messages.

**New plugins:**

- Quotes Plugin: Lets users/viewers or mods add quotes and randomly fetch one.

## [0.0.4](https://git.abyle.org/hps/abylebotter/-/tree/0.0.4) (2020-04-21)

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

- Discord: Initial mapping of AbyleBotter to Discord emojis.
- Discord: Initial converter between AbyleBotter and Discord text format.
- Discord: Support for updating/deleting message.
- Discord: Support for receiving reactions to messages.
- Slack: Support for updating messages.
- Slack: Support for receiving reactions to messages.

**Closed issues:**
- Discord: User.Name shall be User#1234 instead of just User
- Discord: Bot should be able to recover in case of errors
- Discord: Make Callback URL configurable

## [0.0.3](https://git.abyle.org/hps/abylebotter/-/tree/0.0.3) (2020-03-31)

*Major rework of the whole project to a modern and more Go-like structure.*

**Implemented enhancements:**

- Migrated to a server architecture where AbyleBotter is controlled via REST API
- Introduce the command line tool BotterControl to control a AbyleBotter instance
- Cleaner plugin interface and therefore much easier to implement new plugins

**Fixed bugs:**

- Various bugs fixed all over the place

## [0.0.2](https://git.abyle.org/hps/abylebotter/-/tree/0.0.2) (2018-10-06)

*First release.*
