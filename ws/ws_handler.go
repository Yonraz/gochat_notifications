package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/yonraz/gochat_notifications/constants"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub        *redis.Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	Clients    map[string]*Client
}

func NewHandler(h *redis.Client) *Handler {
	return &Handler{
		hub: h,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
		Clients:    make(map[string]*Client),
	}
}

func (h *Handler) JoinNotifications(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("error upgrading to ws connection: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		conn.Close()
		return
	}

	username := ctx.Query("username")
	if username == "" {
		log.Println("Username query parameter missing")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username query parameter is required"})
		conn.Close()
		return
	}
	client := &Client{
		Conn: conn,
		Username: username,
		Message: make(chan *Message),
	}
	message := &Message{
		Sender: username,
		Content: "online",
		Type: constants.UserOnline,
	}

	h.Register <- client
	h.Broadcast <- message

	go client.writePump() // handle writes on other goroutine
	client.readPump(h)
}

func (h *Handler) handleDisconnection(client *Client) {
	message := &Message{
		Sender: client.Username,
		Content: "offline",
		Type: constants.UserOffline,
	}

	h.Broadcast <- message
	fmt.Printf("Client %v left!\n", client.Username)
	h.hub.SRem(context.Background(), "notifications:clients", client.Username)
	h.hub.Del(context.Background(), client.Username)

	delete(h.Clients, client.Username)
}

func (h *Handler) handleConnection(client *Client) {
	fmt.Printf("Client %v joined!\n", client.Username)	
	h.Clients[client.Username] = client
	h.hub.SAdd(context.Background(), "notifications:clients", client.Username)
}

func (h *Handler) Run() {
	for {
		select {
		case client := <-h.Register:
			h.handleConnection(client)
		case client := <-h.Unregister:
			h.handleDisconnection(client)
			client.Conn.Close()
		case message := <-h.Broadcast:
			c := context.Background()
			// Get list of connected clients from Redis
			clients, err := h.hub.SMembers(c, "notifications:clients").Result()
			if err != nil {
				log.Println("Error fetching clients:", err)
				continue
			}

			// Send the message to each client
			for _, username := range clients {
				if client, ok := h.Clients[username]; ok {
					if message.Sender == username {
						continue
					}
					select {
					case client.Message <- message:
						// Successfully sent message to client
					default:
						// Handle case where message channel might be full
						log.Println("Message channel full for client:", username)
					}
				}
			}
		}
	}
}