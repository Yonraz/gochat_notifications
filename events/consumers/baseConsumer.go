package consumers

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_notifications/constants"
	"github.com/yonraz/gochat_notifications/ws"
)

type Consumer struct {
	wsHandler *ws.Handler
	channel     *amqp.Channel
	queueName   string
	routingKey  string
	exchange    string
	handlerFunc func(*ws.Handler, amqp.Delivery) error
}

func NewConsumer(wsHandler *ws.Handler, channel *amqp.Channel, queueName constants.Queues, routingKey constants.RoutingKey, exchange constants.Exchange, handlerFunc func(*ws.Handler, amqp.Delivery) error) *Consumer {
	return &Consumer {
		wsHandler: wsHandler,
		channel: channel,
		queueName: string(queueName),
		routingKey: string(routingKey),
		exchange: string(exchange),
		handlerFunc: handlerFunc,
	}
}

func (c *Consumer) Consume() error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming %w", err)
	}

	go func () {
		for msg := range msgs {
			if err := c.handlerFunc(c.wsHandler, msg); err != nil {
				fmt.Printf("error consuming message %v: %v\n", msg, err)
				msg.Nack(false, false)		
			} else {

			msg.Ack(false)
			}
		} 
	}()

	fmt.Printf("Started consuming on queue: %s\n", c.queueName)
	return nil
}