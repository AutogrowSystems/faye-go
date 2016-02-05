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

## Acknowledgements

* [roncohen](https://github.com/roncohen)

## Licence

```
The MIT License (MIT)

Copyright (c) 2014 Ron Cohen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```