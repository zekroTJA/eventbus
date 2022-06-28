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
