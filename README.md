# Readme

[![Go Report Card](https://goreportcard.com/badge/broadcastle.co/code/emitter-agent)](https://goreportcard.com/report/broadcastle.co/code/emitter-agent) [![api documentation](http://b.repl.ca/v1/api-documentation-blue.png)](https://godoc.org/broadcastle.co/code/emitter-agent)

## Introduction

This project is based off [emitter-actor](https://github.com/kelindar/emitter-actor/).

## Basic Usage

### Requirements

[Emitter](https://github.com/emitter-io/emitter/) is required for this project to work.

Keys must have read and write access. The target channel for the keys must end with `/#/`, ex: `actor/#/`, for this project to work.

### Installation

```bash
go get broadcastle.co/code/emitter-agent
```

### Server

```go
package main

import (
    "fmt"

    agent "broadcastle.co/code/emitter-agent"
)

type data struct {
    URL string `json:"url"`
}

func main() {

    // Don't actually ignore errors.

    logs := logrus.New()

    config := agent.UseEmitter("127.0.0.1", "8090", "key-generate-from-keygen", "actor", logs, false)

    server, _ := agent.New(config)

    dt, _ := agent.PrepSender("server", "ping", data{"https://google.com"})

    server.Send(dt)

}
```

### Client

```go
package main

import (
    "fmt"
    "json"
    "net/http"

    agent "broadcastle.co/code/emitter-agent"
)

type data struct {
    URL string `json:"url"`
}

func main() {

    // Don't actually ignore errors.

    logs := logrus.New()

    config := agent.UseEmitter("127.0.0.1", "8090", "key-generate-from-keygen", "actor", logs, false)

    client, _ := agent.New(config)

    client.Do("ping", func(_ agent.Client, content string) {
    
        var d data
        
        json.Unmarshal([]byte(content), &d)

        res, _ := http.Get(d.URL)

        if res.StatusCode == 200 {
            fmt.Println("success")
        } else {
            fmt.Println("failure")
        }

    })

}
```
