package main
import (
	"fmt"

 	"github.com/RussNavas/learn-pub-sub-starter/internal/pubsub"
	"github.com/RussNavas/learn-pub-sub-starter/internal/gamelogic"
	amqp "github.com/rabbitmq/amqp091-go"
)



func handlerWar(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
    return func(dw gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome, winner, loser:= gs.HandleWar(dw)
		switch outcome{
		
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			err := pubsub.PublishGameLog(
				publishCh,
				gs.Player.Username,
				msg,
			)
			if err != nil{
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			err := pubsub.PublishGameLog(
				publishCh,
				gs.Player.Username,
				msg,
			)
			if err != nil{
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			msg := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			err := pubsub.PublishGameLog(
				publishCh,
				gs.Player.Username,
				msg,
			)
			if err != nil{
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			fmt.Println("outcome error")
			return pubsub.NackDiscard
		}
    }
}
