# Steven

A work in progress Minecraft client in Go.
Don't expect it to go anywhere, just doing this for fun.

## Images

![Steven](http://i.imgur.com/VVnmbkV.png)

![Steven's server list after disconnecting from a server](https://i.imgur.com/JRFXt0e.png)

## Building

To build:

```
export GOPATH=your/install/directory
go get github.com/thinkofdeath/phteven/cmd/phteven
```

To update, run `go get` with the `-u` option.

## What works

* Connecting to servers
* Online mode
* Rendering most blocks
* Block model support

## What doesn't work

* 99% of Minecraft's features

## Chat

I generally am on the `irc.spi.gt` irc network in the `#think` channel. 
Feel free to pop in to say hi, [Webchat can be found here](https://irc.spi.gt/iris/?channels=think)

## Builds

Builds for Linux (64bit only) and Windows(32bit and 64bit) can be found
[Here](http://ci.thinkofdeath.co.uk/viewType.html?buildTypeId=Steven_Client&guest=1)

## Running

![Profile example](http://i.imgur.com/NBMGhPL.png)

You need to create a new profile (or edit an existing one) on the Minecraft 
launcher and modify the profile to look like the above but replace the path
to phteven with the location you built it at or downloaded it too and change the 
`server` parameter to the target server. Currently only works in online mode
(with no plans for offline mode currently).

It is possible to run phteven without the launcher, but you must obtain the access token,
UUID (whithout dashes) and the username, and pass them as arguments to phteven, as well as
the server.
