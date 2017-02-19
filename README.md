# Internet Relay Chat in Go

This module contains my humble attempt of writing an IRC client library.
This is my *very first* Go-based project, so my sources are probably
not what anyone would call "idiomatic Go".

## Developing the Application

Check out the sources into your Go workspace.
```sh
mkdir -p "${GOPATH}/src/github.com/headcr4sh"
git clone https://github.com/headcr4sh/irc.git "${GOPATH}/src/github.com/headcr4sh/irc"
```

### Fetching Build Dependencies
The easiest (and recommended) way to fetch all dependencies is to use "dep": 
```sh
dep ensure
```

### Building the Application
```sh
go generate github.com/headcr4sh/irc/cmd/irc
go install github.com/headcr4sh/irc/cmd/irc
```

## Usage of the Library
The following example shows how to use the event-driven IRC library in your own code,
e.g. to develop a chat bot or something similar.

```go
package main

import (
	"fmt"
	"github.com/headcr4sh/irc"
)

func main() {
	conn := irc.NewClientConnection("irc.freenode.org", 6667)
	go func() {
	    for {
            select {
            case msg := <- conn.In():
                fmt.Printf("Message received: %s", msg)
            case s := <- conn.State():
                fmt.Printf("Connection state changed: %v", s)
            }
		}
	}()
	if err := conn.Open(); err != nil {
		panic("Connection to server failed.")
	}
	conn.Wait()
}

```
