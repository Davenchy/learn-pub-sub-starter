package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(internal.GetRabbitMQURL())
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	defer ch.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		gamelogic.PrintServerHelp()
		for {
			inputs := gamelogic.GetInput()
			for _, input := range inputs {
				switch input {
				case "":
					continue
				case "pause":
					fmt.Println("Pausing game...")
					if err := pubsub.PublishJSON(
						ch, routing.ExchangePerilDirect,
						routing.PauseKey,
						routing.PlayingState{IsPaused: true}); err != nil {
						log.Fatal("Failed to pause the game:", err)
					}
				case "resume":
					fmt.Println("Resuming game...")
					if err := pubsub.PublishJSON(
						ch, routing.ExchangePerilDirect,
						routing.PauseKey,
						routing.PlayingState{IsPaused: false}); err != nil {
						log.Fatal("Failed to resume the game:", err)
					}
				case "quit":
					cancel()
					return
				default:
					fmt.Println("Unknown command:", input)
				}
			}
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down Peril server...")
}
