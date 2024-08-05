package ws

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yonraz/gochat_notifications/constants"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Conn     *websocket.Conn
	Message  chan *Message
	Username string
	IsGuest bool 
}

type Message struct {
	ID string `json:"ID" gorm:"primarykey"`
	Sender string `json:"sender"`
	Receiver string `json:"receiver"`
	Content string `json:"content"`
	Type constants.Notification `json:"type"`
	CreatedAt time.Time  				`json:"createdAt"`
	UpdatedAt time.Time  				`json:"updatedAt"`
}

func (client *Client) readPump(handler *Handler) {
	defer func() {
		handler.handleDisconnection(client)
		// handler.handleDisconnection(client)
		client.Conn.Close()
	}()
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("connection closed: %v\n", err)
			}
			break
		}

		handler.Broadcast <- &Message{
			Content: string(message),
			Sender: client.Username,
		}
	}
}

func (client *Client) writePump() {
	defer func() {
			client.Conn.Close()
	}()
	for message := range client.Message {
		fmt.Printf("Sending message %v\n", message)
		client.Conn.WriteJSON(message)
	}
}