package utils

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_notifications/constants"
)

type queueConstructor struct {
	Queue constants.Queues
	Key constants.RoutingKey
	Exchange constants.Exchange
}

func DeclareAndBindQueue(
	channel *amqp.Channel,
	queueName constants.Queues,
	routingKey constants.RoutingKey,
	exchangeName constants.Exchange,
	durable bool,
	autoDelete bool,
	exclusive bool,
	args amqp.Table,
) error {
	_, err := channel.QueueDeclare(
		string(queueName),
		durable,
		autoDelete,
		exclusive,
		false,
		args,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = channel.QueueBind(
		string(queueName),
		string(routingKey),
		string(exchangeName),
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

func DeclareQueues(channel *amqp.Channel) error {
	queues := []queueConstructor{
		{Queue: constants.MessageDeliveredQueue, Key: constants.MessageDeliveredKey, Exchange: constants.MessageEventsExchange},
		{Queue: constants.MessageSentQueue, Key: constants.MessageSentKey, Exchange: constants.MessageEventsExchange},
		{Queue: constants.MessageReadQueue, Key: constants.MessageReadKey, Exchange: constants.MessageEventsExchange},
	}

	for _, q := range queues {
		err := DeclareAndBindQueue(
			channel,
			q.Queue,
			q.Key,
			q.Exchange,
			true,
			false,
			false,
			nil,
		)
		if err != nil {
		return err
		}
	}

	
	return nil
}