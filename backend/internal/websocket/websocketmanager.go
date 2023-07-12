package websocket

import (
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/models"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // or check the origin if you want to add more security
	},
}

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

type ClientMessage struct {
	Client  *Client
	Message []byte
}

type WebSocketManager struct {
	Clients     map[*Client]bool
	ClientsByID map[string]*Client
	Broadcast   chan ClientMessage
	Register    chan *Client
	Unregister  chan *Client
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Broadcast:   make(chan ClientMessage),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[*Client]bool),
		ClientsByID: make(map[string]*Client),
	}
}

func (manager *WebSocketManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.Clients[client] = true
			manager.ClientsByID[client.UserID] = client
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client]; ok {
				delete(manager.Clients, client)
				delete(manager.ClientsByID, client.UserID)
				close(client.Send)
			}
		case clientMessage := <-manager.Broadcast:
			for client := range manager.Clients {
				if client == clientMessage.Client {
					// Skip the client that sent the message
					continue
				}
				select {
				case client.Send <- clientMessage.Message:
				default:
					close(client.Send)
					delete(manager.Clients, client)
					delete(manager.ClientsByID, client.UserID)
				}
			}
		}
	}
}

func (manager *WebSocketManager) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	// Read the first message from the client
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.WithError(err).Error("Error reading message")
		conn.Close()
		return
	}

	// Unmarshal the message into a map
	var msgMap models.ConnectionSetup
	if err := proto.Unmarshal(message, &msgMap); err != nil {
		log.WithError(err).Error("Error unmarshalling message")
		conn.Close()
		return
	}

	// Check that the message type is "register"
	if msgMap.Type != "register" {
		log.Error("First message type is not 'register'")
		conn.Close()
		return
	}

	// Get the user ID from the message
	userID := msgMap.UserID
	if userID == "" {
		log.Error("No userID field in register message")
		conn.Close()
		return
	}

	client := &Client{Conn: conn, Send: make(chan []byte)}

	manager.Register <- client

	go manager.HandleSending(client)
	manager.HandleReceiving(client)
}

func (manager *WebSocketManager) HandleSending(client *Client) {
	defer func() {
		manager.Unregister <- client
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			client.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (manager *WebSocketManager) HandleReceiving(client *Client) {
	defer func() {
		manager.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			manager.Unregister <- client
			client.Conn.Close()
			break
		}

		manager.Broadcast <- ClientMessage{Client: client, Message: message}
	}
}

func (manager *WebSocketManager) SendMessageToClient(client *Client, message []byte) {
	// Ensure the client is currently connected
	if _, ok := manager.Clients[client]; ok {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(manager.Clients, client)
		}
	}
}

func (manager *WebSocketManager) SendMessageToUserID(userID string, message []byte) {
	client, ok := manager.ClientsByID[userID]
	if ok {
		manager.SendMessageToClient(client, message)
	}
}
