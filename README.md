<div align="center">
    <h1>~ eventbus ~</h1>
    <strong>A go package to send and receive pub-sub messages using channels.</strong><br><br>
    <a href="https://pkg.go.dev/github.com/zekrotja/eventbus"><img src="https://godoc.org/github.com/zekrotja/eventbus?status.svg" /></a>&nbsp;
    <a href="https://github.com/zekrotja/timedmap/actions/workflows/tests.yml" ><img src="https://github.com/zekroTJA/timedmap/actions/workflows/tests.yml/badge.svg" /></a>&nbsp;
    <a href="https://coveralls.io/github/zekrotja/eventbus"><img src="https://coveralls.io/repos/github/zekrotja/eventbus/badge.svg" /></a>&nbsp;
    <a href="https://goreportcard.com/report/github.com/zekrotja/eventbus"><img src="https://goreportcard.com/badge/github.com/zekrotja/eventbus"/></a>
<br>
</div>

---

<div align="center">
    <code>go get -u github.com/zekrotja/eventbus</code>
</div>

---

## Intro

This package provides a very simple, generic pub-sub event bus to simplify sending event 
messages between services using channels.

### Why using channels over callbacks?

Callback are - in my opinion - not a very clean and performant way to perform event driven
development in Go because, in contrast to languages like JavaScript, Go is not an event
driven language. This can lead to blocking publisher routines while waiting for the execution
of callbacks on the side of subscribers. Channels, which are well designed to communicate
between go routines in the first place, are therefore a way better tool to achieve easy, performant
and intuitive communication of events between publishers and subscribers.

## Basic Example

```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/zekrotja/eventbus"
)

func subscriber(name string, bus *eventbus.EventBus[string]) {
	c, _ := bus.Subscribe()
	for msg := range c {
		fmt.Printf("[ -> %s ]: %s\n", name, msg)
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	bus := eventbus.New[string]()

	go subscriber("service1", bus)
	go subscriber("service2", bus)

	bus.SubscribeFunc(func(s string) {
		if s == "exit" {
			os.Exit(0)
		}
	})

	fmt.Println("Publish messages by writing them to the console.\nPress CTRL+C or write 'exit' to exit.")
	for s.Scan() {
		bus.Publish(s.Text())
	}
}
```

Further examples can be found in the [examples](examples/) directory. If you want to take alook at a
practical example, feel free to explore my project [Yuri69](https://github.com/zekrotja/yuri69) which
heavily depends on EventBus.