package main

import (
	"fmt"

	"github.com/RussNavas/learn-pub-sub-starter/internal/gamelogic"
	"github.com/RussNavas/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState){
	return func(r routing.PlayingState){
		defer fmt.Print("> ")
		gs.HandlePause(r)
	}
}
