package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yonraz/gochat_notifications/controllers"
	"github.com/yonraz/gochat_notifications/events/consumers"
	"github.com/yonraz/gochat_notifications/initializers"
	"github.com/yonraz/gochat_notifications/ws"
)

func init() {
	fmt.Println("Application starting...")
	time.Sleep(1 * time.Minute)
	initializers.LoadEnvVariables()
	initializers.ConnectToRabbitmq()
	initializers.ConnectToRedis()
}

func main() {
	router:= gin.Default()
	defer func() {
		if err := initializers.RmqChannel.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		}
	}()
	defer func() {
		if err := initializers.RmqConn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		}
	}()
	wsHandler := ws.NewHandler(initializers.RedisClient)
	notifController := controllers.NewNotificationsController(initializers.RedisClient)

	router.GET("/ws/notifications/:username", notifController.GetNotificationsForUser)
	router.DELETE("ws/notifications/:username", notifController.DeleteNotifications)
	router.GET("/ws/notifications/join", wsHandler.Join)
	
	go wsHandler.Run()
	
	messageSentConsumer := consumers.NewMessageSentConsumer(wsHandler, initializers.RmqChannel)

	go messageSentConsumer.Consume()
	router.Run()
}