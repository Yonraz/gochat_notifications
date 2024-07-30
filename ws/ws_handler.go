package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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
	sync.RWMutex
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

func (h *Handler) Join(ctx *gin.Context) {
	username := ctx.Query("username")
	h.Lock()
	if _, exists := h.Clients[username]; exists {
		log.Printf("Username %v already connected", username)
		ctx.JSON(http.StatusConflict, gin.H{"error": "User already connected"})
		h.Unlock()
		return
	}
	h.Unlock()
	
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("error upgrading to ws connection: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		conn.Close()
		return
	}
	isGuest := username == ""
	if isGuest {
		username = uuid.New().String()
	} 
	client := &Client{
		Conn: conn,
		Username: username,
		Message: make(chan *Message),
		IsGuest: isGuest,
	}
	h.Register <- client
	
	if !isGuest {
		message := &Message{
			Sender: username,
			Content: "online",
			Type: constants.UserOnline,
		}
		h.Broadcast <- message
	}
	
	log.Printf("initializing connection for %v...\n", username)
	go client.writePump() 
	go client.readPump(h)
}


func (h *Handler) handleDisconnection(client *Client) {
	
	h.Lock()
	defer h.Unlock()
	if _, ok := h.Clients[client.Username]; ok {
		client.Conn.Close()
		delete(h.Clients, client.Username)
		if !client.IsGuest {
			message := &Message{
				Sender: client.Username,
				Content: "offline",
				Type: constants.UserOffline,
			}

			h.Broadcast <- message
		}
		
		fmt.Printf("Client %v left!\n", client.Username)
	}
	

	// Log initial state
	log.Printf("Disconnecting client: %s", client.Username)



	// Remove client from Redis set
	removeCmd := h.hub.SRem(context.Background(), string(constants.NotificationClients), client.Username)
	if err := removeCmd.Err(); err != nil {
		log.Printf("Error removing client from Redis set: %v", err)
	} else if removedCount := removeCmd.Val(); removedCount == 0 {
		log.Printf("Client %s was not found in Redis set", client.Username)
	} else {
		log.Printf("Removed client %s from Redis set", client.Username)
	}
	log.Printf("Final Redis clients: %v", h.hub.SMembers(context.Background(), string(constants.NotificationClients)).Val())

}

func (h *Handler) handleConnection(client *Client) {
	fmt.Printf("Client %v joined!\n", client.Username)	

	h.Lock()
	defer h.Unlock()
	h.Clients[client.Username] = client
	h.hub.SAdd(context.Background(), string(constants.NotificationClients), client.Username)
}
func (h *Handler) handleBroadcast(message *Message) {
    c := context.Background()
    clients, err := h.hub.SMembers(c, "notifications:clients").Result()
    if err != nil {
        log.Println("Error fetching clients:", err)
        return
    }
    fmt.Printf("%v", clients)

    for _, username := range clients {
        h.RLock()
        client, ok := h.Clients[username]
        h.RUnlock()
        if !ok || message.Sender == username {
            continue
        }
        select {
        case client.Message <- message:
        default:
            log.Println("Message channel full for client:", username)
        }
    }
}

func (h *Handler) Run() {
	for {
		select {
		case client := <-h.Register:
			log.Println("in register")
			h.handleConnection(client)
		case client := <-h.Unregister:
			log.Println("in unregister")
			log.Printf("client for unregister is %v", client)
			h.handleDisconnection(client)
			client.Conn.Close()
		case message := <-h.Broadcast:
			h.handleBroadcast(message)
		}
	}
}