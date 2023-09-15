package websocket

import (
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"net/http"
)

var WsManager = NewWebSocketManager()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement a proper check for security reasons.
		return true
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
		}
	}
}

func (manager *WebSocketManager) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.WithError(err).Error("Failed to upgrade the connection")
		return
	}

	// Read the first message from the client
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.WithError(err).Error("Error reading message")
		conn.Close()
		return
	}

	// Unmarshal the message into a map
	var msgMap marketplace.ConnectionSetup
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

	log.Info("New Client")
	client := &Client{Conn: conn, Send: make(chan []byte, 10000), UserID: msgMap.UserID}

	log.Info("Registering client")
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
				err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.WithError(err).Error("Error writing message to websocket")
				}
				return
			}

			err := client.Conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.WithError(err).Error("Error writing message to websocket")
			}

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
	manager.sendMessage(client, message)
}

func (manager *WebSocketManager) SendMessageToAllClients(message []byte) {
	for client := range manager.Clients {
		manager.sendMessage(client, message)
	}
}

func (manager *WebSocketManager) SendMessageToUserID(userID string, message []byte) {
	if client, ok := manager.ClientsByID[userID]; ok {
		manager.sendMessage(client, message)
	}
}

func (manager *WebSocketManager) isConnected(client *Client) bool {
	_, ok := manager.Clients[client]
	return ok
}

func (manager *WebSocketManager) sendMessage(client *Client, message []byte) {
	if !manager.isConnected(client) {
		log.Error("Client not connected")
		return
	}

	select {
	case client.Send <- message:
	default:
		log.Error("Closing socket block sender")
		close(client.Send)
		delete(manager.Clients, client)
	}
}
