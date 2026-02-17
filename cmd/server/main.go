package main

import (
	"fmt"

	"github.com/RussNavas/learn-pub-sub-starter/internal/gamelogic"
	"github.com/RussNavas/learn-pub-sub-starter/internal/pubsub"
	"github.com/RussNavas/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil{
		fmt.Printf("problem creating *connection: %v", err)
		return
	}
	defer connection.Close()
	fmt.Println("Connection Successful!")

	ch, err := connection.Channel()
	if err != nil{
		fmt.Printf("problem creating *channel %v", err)
		return
	}

	defer ch.Close()

	gamelogic.PrintServerHelp()
	
	for {

		input := gamelogic.GetInput()
		if len(input) == 0{
			continue
		}

		firstWord := input[0]
		switch firstWord{

			case "pause":
				fmt.Println("The game has been paused")
				err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
					IsPaused: true,
				})
				if err != nil{
					fmt.Printf("problem publishing json")
					return
				}

			case "resume":
				fmt.Println("The game has been resumed.")
				err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
					IsPaused: false,
				})
				if err != nil{
					fmt.Printf("problem publishing json")
					return
				}
			case "quit":
				fmt.Println("Exiting the game")
				return
			default:
				fmt.Println("invalid command try again...")
				continue
		}


	}

	/*
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Printf("Shutting down & closing connection")
	*/
}
