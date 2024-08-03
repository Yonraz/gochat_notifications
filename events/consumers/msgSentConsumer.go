package consumers

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_notifications/constants"
	"github.com/yonraz/gochat_notifications/models"
	"github.com/yonraz/gochat_notifications/ws"
)

func NewMessageSentConsumer(wsHandler *ws.Handler, channel *amqp.Channel) *Consumer {
	return &Consumer{
		wsHandler: wsHandler,
		channel: channel,
		queueName: string(constants.MessageSentQueue),
		routingKey: string(constants.MessageSentKey),
		exchange: string(constants.MessageEventsExchange),
		handlerFunc: MessageSentHanlder,
	}
}

func MessageSentHanlder(wsHandler *ws.Handler, msg amqp.Delivery) error {
	var parsed *models.WsMessage
	
	json.Unmarshal(msg.Body, &parsed)
	fmt.Printf("message %v consumed on exchange %v with routing key %v\n", parsed, constants.MessageEventsExchange, constants.MessageSentKey)
	message := &ws.Message{
		ID: parsed.ID,
		Content: fmt.Sprintf("Message received from %v", parsed.Sender),
		Sender: parsed.Sender,
		Type: parsed.Type,
		Receiver: parsed.Receiver,
	}
	wsHandler.Broadcast <- message
	return nil
}