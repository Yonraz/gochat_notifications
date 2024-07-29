package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

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
	if username == "" {
		log.Println("Username query parameter missing")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username query parameter is required"})
		return
	}
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
	log.Printf("initializing connection for %v...\n", username)

	h.Register <- client
	h.Broadcast <- message

	go client.writePump() 
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

	h.Lock()
	defer h.Unlock()
	h.hub.SRem(context.Background(), string(constants.NotificationClients), client.Username)
	h.hub.Del(context.Background(), client.Username)

	delete(h.Clients, client.Username)
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