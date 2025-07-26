package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bootdotdev/learn-pub-sub-starter/internal"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial(internal.GetRabbitMQURL())
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("Failed to get username: ", err)
	}

	ch, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.TransientQType)
	if err != nil {
		log.Fatal("Failed to declare and bind queue:", err)
	}
	defer ch.Close()

	game := gamelogic.NewGameState(username)
	// REPL
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			fmt.Println()
			continue
		}

		switch words[0] {
		case "spawn":
			if err := game.CommandSpawn(words); err != nil {
				log.Println("Failed to spawn a unit:", err)
			}
		case "move":
			if _, err := game.CommandMove(words); err != nil {
				log.Println("Error:", err)
			}
		case "spam":
			log.Println("Spamming not allowed yet!")
		case "status":
			game.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			fmt.Println("Shutting down Peril client...")
			os.Exit(0)
		default:
			log.Println("Invalid command. Try again:", words[0])
		}
	}
}
