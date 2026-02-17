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

	isDurable := false
	if queueType == SimpleQueueDurable{
		isDurable = true
	}

	isAutoDelete := false
	isExclusive := false
	if queueType == SimpleQueueTransient{
		isAutoDelete = true
		isExclusive = true
	}

	q, err := channel.QueueDeclare(
		queueName,
		isDurable,
		isAutoDelete,
		isExclusive,
		false,
		nil)

	if err != nil{
		return channel, q, fmt.Errorf("problem with QueueDeclare: %v", err)
	}
	
	err = channel.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil)

	if err != nil{
		return channel, q, fmt.Errorf("problem binding queue: %v", err)
	}

	return channel, q, nil
}

