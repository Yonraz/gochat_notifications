package utils

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_notifications/constants"
)

func DeclareExchanges(channel *amqp.Channel) error {
	err := channel.ExchangeDeclare(
		string(constants.UserEventsExchange),
		"topic",             
		true,                
		false,               
		false, 
		false,              
		nil,   
	)
	if err != nil {
		return err
	}

	err = channel.ExchangeDeclare(
		string(constants.MessageEventsExchange),
		"topic",             
		true,                
		false,               
		false, 
		false,              
		nil,   
	)
	if err != nil {
		return err
	}
	fmt.Printf("Exchanges %v created!\n", constants.MessageEventsExchange)
	return nil
}