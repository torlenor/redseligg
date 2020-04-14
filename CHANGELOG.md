# Change Log

## [0.0.4](https://git.abyle.org/hps/abylebotter/-/tree/0.0.4) (TBD)

**Implemented enhancements:**

- Added GetVersion() to the Bot/Plugin API.
- Added OnRun() hook from Bot to Plugin which is called when the Bot is ready to serve the plugins.
- Added OnStop() hook from Bot to Plugin which is called when the Bot is shutting down.
- Added OnReactionAdded, OnReactionRemoved hooks from Bot to Plugin.

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
