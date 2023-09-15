package main

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("URL must be provided as the first argument.")
	}
	rawURL := os.Args[1]

	// Validate URL
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}
	// Create a request to pass along the basic auth header
	requestHeader := make(http.Header)
	requestHeader.Set("Authorization", "Basic "+basicAuth("admin", "unity"))

	// Connect to WebSocket with Basic Authentication
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), requestHeader)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()
	defer conn.Close()

	if err := sendProtobufMessage(conn, generateConnection(), deployEks()); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	poll()
	fmt.Println("Message sent successfully.")
}

func poll() {

}

func generateConnection() *marketplace.ConnectionSetup {
	return &marketplace.ConnectionSetup{
		Type:   "register",
		UserID: "integrationtest",
	}
}

func deployEks() *marketplace.UnityWebsocketMessage {
	app := &marketplace.Install_Applications{
		Name:        "unity-eks",
		Version:     "0.1",
		Variables:   nil,
		Postinstall: "",
		Preinstall:  "",
	}
	inst := &marketplace.Install{
		Applications:   app,
		DeploymentName: "nightly-eks",
	}

	instm := &marketplace.UnityWebsocketMessage_Install{Install: inst}
	return &marketplace.UnityWebsocketMessage{Content: instm}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func sendProtobufMessage(conn *websocket.Conn, connect *marketplace.ConnectionSetup, msg *marketplace.UnityWebsocketMessage) error {
	if connect != nil {
		data, err := proto.Marshal(connect)
		if err != nil {
			return fmt.Errorf("failed to marshal protobuf: %v", err)
		}

		return conn.WriteMessage(websocket.BinaryMessage, data)
	}

	if msg != nil {
		data, err := proto.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal protobuf: %v", err)
		}

		return conn.WriteMessage(websocket.BinaryMessage, data)
	}
	return nil
}
