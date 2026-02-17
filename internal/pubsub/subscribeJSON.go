package pubsub

import (
	"fmt"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T),
) error {

	channel, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)

	if err != nil {
		return fmt.Errorf("DeclarAndBind failed: %v, Queue: %v", err, queue.Name)
	}

	msgs, err := channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,)
	if err != nil{
		return fmt.Errorf("unable to use channel.Consume: %v", err)
	}

	go func(){
		for msg := range msgs {
			var val T
			err := json.Unmarshal(msg.Body, &val)
			if err != nil{
				fmt.Errorf("unable to unmarshal body: %v, channel: %v", err, msg)
			}
			handler(val)
			err = msg.Ack(true)
			if err != nil{
				fmt.Errorf("unable to Ack: %v", err)
			}
		}
	}()
	return nil

}


