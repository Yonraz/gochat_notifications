package consumers

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_notifications/constants"
	"github.com/yonraz/gochat_notifications/models"
)

func NewMessageSentConsumer(channel *amqp.Channel) *Consumer {
	return &Consumer{
		channel: channel,
		queueName: string(constants.MessageSentQueue),
		routingKey: string(constants.MessageSentKey),
		exchange: string(constants.MessageEventsExchange),
		handlerFunc: MessageSentHanlder,
	}
}

func MessageSentHanlder(msg amqp.Delivery) error {
	var parsed *models.WsMessage
	
	json.Unmarshal(msg.Body, &parsed)
	fmt.Printf("message %v consumed on exchange %v with routing key %v\n", parsed, constants.MessageEventsExchange, constants.MessageSentKey)
	return nil
}