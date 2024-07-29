package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
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
	router.GET("/ws/notifications/join", wsHandler.Join)
	
	go wsHandler.Run()
	router.Run()
}