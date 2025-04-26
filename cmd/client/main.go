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
	url := "amqp://guest:guest@localhost:5672"
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection successful")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// _, err = conn.Channel()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println(err)
		return
	}
	ch, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, "pause."+username, routing.PauseKey, 1)
	fmt.Println(ch, queue)
	if err != nil {
		fmt.Println(err)
		return
	}
	if <-c == os.Interrupt {
		fmt.Println("\n>Shutting down")
		return
	}
}
