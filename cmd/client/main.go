package main

import (
	"fmt"
	"log"

	"github.com/RussNavas/learn-pub-sub-starter/internal/gamelogic"
	"github.com/RussNavas/learn-pub-sub-starter/internal/pubsub"
	"github.com/RussNavas/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil{
		fmt.Printf("problem creating *connection: %v", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection Successful!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}


	gs := gamelogic.NewGameState(username)

	pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		"pause."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gs),
	)

	chann, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		"army_moves."+username,
		"army_moves.*",
		pubsub.SimpleQueueTransient,
	)
	if err != nil{
		log.Fatalf("couldnt DeclareAndBind army_moves: %v", err)
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "move":
			mv, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}

			gs.HandleMove(mv)
			defer fmt.Print("> ")

			err = pubsub.PublishJSON(
				chann,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+username,
				mv, 
			)
			if err != nil {
				fmt.Print(fmt.Errorf("problem publishing move to topic exchange: %w", err))
			}
			fmt.Printf("Move Succesful\n")


		case "spawn":
			err = gs.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			// TODO: publish n malicious logs
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}

}
