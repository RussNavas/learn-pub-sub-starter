package main

import (
	"os"
	"os/signal"
	"fmt"
	"github.com/RussNavas/learn-pub-sub-starter/internal/gamelogic"
	"github.com/RussNavas/learn-pub-sub-starter/internal/pubsub"
	"github.com/RussNavas/learn-pub-sub-starter/internal/routing"
    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil{
		fmt.Printf("problem creating *connection: %v", err)
		return
	}
	defer connection.Close()
	fmt.Println("Connection Successful!")

	userName, err := gamelogic.ClientWelcome()
	if err != nil{
		fmt.Printf("Problem getting username: %v", err)
		return
	}

	_, _, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		routing.PauseKey +"."+userName,
		routing.PauseKey,
		pubsub.SimpleQueueTransient)

	if err != nil{
		fmt.Printf("Problem with DeclareAndBind: %v", err)
		return
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Printf("Shutting down & closing connection")


}
