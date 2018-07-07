# AbyleBotter

## Description

This is AbyleBotter, an extensible Chat Bot for Discord (and in the future other chat platforms).

At the moment the Bot is in a proof of concept/API/interface development phase with very limited functional use.

## How to download/install

Make sure your $GOPATH is set.

Checkout AbyleEDA with

    go get github.com/torlenor/AbyleBotter

Checkout all dependencies

    go get -v ./...

Build it with

    go install github.com/torlenor/AbyleBotter

The binaries will land in $GOPATH/bin

## How to start

Currently it is necessary to have a registered Bot Account which has already joined some servers. Please take a look at https://discordapp.com/developers/docs/intro on how to set this up. In the end you should have a bot token to use with AbyleBotter.

To start AbyleBotter use when using Linux (bash or similar shell)

```
DISCORD_BOT_TOKEN=token $GOPATH/bin/AbyleBotter
```

or on Windows open cmd, change into your $GOPATH\bin and type

```
set DISCORD_BOT_TOKEN=token
AbyleBotter.exe
```

The Bot should now connect to the Discord Gateway and should be ready to use. 

## Using Docker

If you want try out AbyleBotter in a Docker container you can use the provided Dockerfile to build a Docker image using

```
docker build -t abylebotter .
```

After the successful build of the docker image, it can be started using

```
docker run -i -t --name abylebotter --rm -e DISCORD_BOT_TOKEN=token abylebotter
```

where token has to be replaced with your Discord Bot token.