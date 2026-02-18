package main

import (
	"fmt"

	"github.com/RussNavas/learn-pub-sub-starter/internal/pubsub"
	"github.com/RussNavas/learn-pub-sub-starter/internal/gamelogic"
	"github.com/RussNavas/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype{
	return func(r routing.PlayingState) pubsub.Acktype{
		defer fmt.Print("> ")
		gs.HandlePause(r)
		return pubsub.Ack
	}
}
