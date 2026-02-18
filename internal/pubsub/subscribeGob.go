package pubsub

import (
	"bytes"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)


func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
) error {

	gobUnmarshaller := func(data []byte) (T, error) {
		var target T
		buf := bytes.NewBuffer(data)
		decoder := gob.NewDecoder(buf)
		err := decoder.Decode(&target)
		return target, err
	}

	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
		handler,
		gobUnmarshaller,
	)
}


func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshaller func([]byte) (T, error),
) error {

	ch, _, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
	)
	if err != nil {
		return fmt.Errorf("could not declare and bind: %v", err)
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("could not consume: %v", err)
	}

	go func() {
		defer ch.Close()
		for msg := range msgs {
			target, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				msg.Nack(false, false) 
				continue
			}

			switch handler(target) {
			case Ack:
				msg.Ack(false)
			case NackDiscard:
				msg.Nack(false, false)
			case NackRequeue:
				msg.Nack(false, true)
			}
		}
	}()

	return nil
}
