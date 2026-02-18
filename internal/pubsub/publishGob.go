package pubsub

import (
	"fmt"
	"bytes"
	"context"
	"encoding/gob"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(val)
	if err != nil{
		return fmt.Errorf("problem with encoding val to gob: %v", err)
	}

	err = ch.PublishWithContext(context.Background(),
			exchange,
			key,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/gob",
				Body: b.Bytes(),
			},
	)
	if err != nil{
		return fmt.Errorf("problem publishing gob with context: %v", err)
	}
	return nil
}


