# Steven

A work in progress Minecraft client in Go.
Don't expect it to go anywhere, just doing this for fun.

## Images

![Steven](http://i.imgur.com/VVnmbkV.png)

![Steven with Nether style background and a resource pack](https://i.imgur.com/QjBb1UT.png)

![Steven's server list after disconnecting from a server](https://i.imgur.com/JRFXt0e.png)

## Building

To build:

```
export GOPATH=your/install/directory
go get github.com/thinkofdeath/steven/cmd/steven
```

To update, run `go get` with the `-u` option.

Requires `csfml` libraries and headers to build. To include these in your build, 
create the following environment variables:

* `CGO_CFLAGS=-Ipath/to/csfml/include`
* `CGO_LDFLAGS=-Lpath/to/csfml/lib`

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

**Latest:**

|  #  |      Linux      | OS X |      Windows      |
|:---:|:---------------:|:----:|:-----------------:|
| x64 | [linux_amd64.zip](http://ci.thinkofdeath.uk/guestAuth/repository/download/Steven_Client/.lastSuccessful/linux_amd64.zip) |   [Issue](https://github.com/thinkofdeath/steven/issues/27)  | [windows_amd64.zip](http://ci.thinkofdeath.uk/guestAuth/repository/download/Steven_Client/.lastSuccessful/windows_amd64.zip) |
| x32 |   [Issue](https://github.com/thinkofdeath/steven/issues/28)       |   [Issue](https://github.com/thinkofdeath/steven/issues/27)  |  [windows_386.zip](http://ci.thinkofdeath.uk/guestAuth/repository/download/Steven_Client/.lastSuccessful/windows_386.zip)  |

Older builds can be found [here](http://ci.thinkofdeath.co.uk/viewType.html?buildTypeId=Steven_Client&guest=1)

## Running

### Via the Offical Minecraft launcher

![Profile example](http://i.imgur.com/NBMGhPL.png)

You need to create a new profile (or edit an existing one) on the Minecraft 
launcher and modify the profile to look like the above but replace the path
to steven with the location you built it at or downloaded it too and change the 
`server` parameter to the target server. Currently only works in online mode
(with no plans for offline mode currently). If the `server` parameter isn't
passed then a server list will be displayed.

### Standalone

Just running steven via a double click (Windows) or `./steven` (everything else)
will bring up a login screen followed by a server list which you can select a server
from.

Providing a username, uuid and access token via the command line as followed:
  `--username <username> --uuid <uuid> --accessToken <access token>`
will skip the login screen and jump straight to the server list. Providing a
server address via `--server <server>:<port>` will skip the server list and 
connect straight to the server. As it currently stands providing all the arguments
allows for the client to parallelise connecting to the server and loading the 
textures/models/other assets as a 'quick connect'.

