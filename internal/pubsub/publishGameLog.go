package pubsub

import (

	"time"

	"github.com/RussNavas/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)


func PublishGameLog(publishCh *amqp.Channel, username, msg string) error {
	return PublishGob(
		publishCh,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+username,
		routing.GameLog{
			CurrentTime: time.Now(),
			Message:	msg,
			Username:	username,
		},
	)
}

