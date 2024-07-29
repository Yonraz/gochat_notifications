package initializers

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_notifications/events/utils"
)

var RmqChannel *amqp.Channel
var RmqConn *amqp.Connection

func ConnectToRabbitmq() {
	var err error
	user := os.Getenv("RMQ_USER")
	password := os.Getenv("RMQ_PASSWORD")
	connectionString := fmt.Sprintf("amqp://%v:%v@rabbitmq:5672/", user, password)
	RmqConn, err = amqp.Dial(connectionString)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Connected to Rabbitmq!")

	RmqChannel, err = RmqConn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = utils.DeclareExchanges(RmqChannel)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = utils.DeclareQueues(RmqChannel)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}