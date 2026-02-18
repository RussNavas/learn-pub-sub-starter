package pubsub

import(
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int
const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
)(*amqp.Channel, amqp.Queue, error){

	channel, err := conn.Channel()
	if err != nil{
		return nil, amqp.Queue{}, fmt.Errorf("problem creating channel on amqp connection: %v", err)
		
	}

	queue, err := channel.QueueDeclare(
		queueName,                       // name
		queueType == SimpleQueueDurable, // durable
		queueType != SimpleQueueDurable, // delete when unused
		queueType != SimpleQueueDurable, // exclusive
		false,                           // no-wait
		amqp.Table{"x-dead-letter-exchange": "peril_dlx"},
	)
	if err != nil{
		return channel, queue, fmt.Errorf("problem with QueueDeclare: %v", err)
	}
	
	err = channel.QueueBind(
		queue.Name,
		key,
		exchange,
		false,
		nil)

	if err != nil{
		return channel, queue, fmt.Errorf("problem binding queue: %v", err)
	}

	return channel, queue, nil
}

