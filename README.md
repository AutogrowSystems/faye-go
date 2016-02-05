# Faye server in Go

This is a Faye server written it in go.  It is still experimentl at this stage.

## Build

    $ go get github.com/AutogrowSystems/faye-go
    $ GO15VENDOREXPIMENT=1 go build github.com/AutogrowSystems/faye-go/main/faye_server.go

## Run

Here's the usage for it:

    Usage of faye_server:
      -h string
        	Port number to serve on (default "127.0.0.1")
      -p string
        	Port number to serve on (default "8000")
      -public string
        	Port number to serve on (default "src/github.com/AutogrowSystems/faye-go/public")

Simply run it like so:

    $ ./faye_server -p 8099

## Features

* websockets
* long-polling