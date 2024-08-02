package ws

import (
	"fmt"
	"log"

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
	gorm.Model
	ID string `json:"ID" gorm:"primarykey"`
	Sender string `json:"sender"`
	Content string `json:"content"`
	Type constants.MessageType `json:"type"`
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