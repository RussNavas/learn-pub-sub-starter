package main

import (
	"fmt"
	"os"
	"os/signal"
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

	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	if err != nil{
		fmt.Printf("problem publishing json")
		return
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Printf("Shutting down & closing connection")
}
