package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	url := "amqp://guest:guest@localhost:5672"
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection successful")

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		words := gamelogic.GetInput()
		if words[0] == "pause" {
			fmt.Println("pausing")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			// IsPaused = true
		} else if words[0] == "resume" {
			fmt.Println("resuming")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			// IsPaused = false
		} else if words[0] == "quit" {
			fmt.Println("quitting")
			break
		} else {
			fmt.Println("i don't understand the cmd")
		}
	}

	pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
	if <-c == os.Interrupt {
		fmt.Println("\n>Shutting down")
		return
	}
}
