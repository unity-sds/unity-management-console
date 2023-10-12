package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/processes"
	websocket2 "github.com/unity-sds/unity-management-console/backend/internal/websocket"
	"net/http"
	"net/http/pprof"
)

var conf config.AppConfig

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan websocket2.ClientMessage)

var appConf config.AppConfig

// setupFeatureFlags sets up feature flags for the application.
// It uses the username from the gin context to create a new user for the feature flag client.
// It then checks the value of the "test-flag" for the user and logs the result.
func setupFeatureFlags(c *gin.Context) {
	log.Info("Setting up feature flags")

	username := c.MustGet(gin.AuthUserKey).(string)
	user := ffuser.NewUser(username)

	hasFlag, _ := config.FFClient.BoolVariation("test-flag", user, false)
	if hasFlag { // flag "test-flag" is true for the user
		log.Info("Flag true")
	} else { // flag "test-flag" is false for the user
		log.Info("flag false")
	}
}

// handleRoot redirects the root URL to "/ui".
func handleRoot(c *gin.Context) {
	log.Infof("Inside handleRoot")
	//c.Redirect(http.StatusMovedPermanently, appConf.BasePath+"/ui")
	c.JSON(http.StatusOK, gin.H{
		"error": "you hit the root",
	})
}

// handlePing responds with a JSON message containing "pong".
// This can be used to check if the server is running.
func handlePing(c *gin.Context) {
	log.Info("Inside handlePing")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// handleWebsocket handles websocket connections.
// It upgrades the HTTP connection to a websocket connection and reads messages from the client.
// Each message is unmarshalled into a WebsocketMessage and sent to the broadcast channel.
func handleWebsocket(c *gin.Context) {
	websocket2.WsManager.HandleConnections(c.Writer, c.Request)
}

// handleNoRoute serves the index.html file for any routes that are not defined.
func handleNoRoute(c *gin.Context) {
	log.Infof("Inside handleNoRoute, mismatched path: %s", c.Request.URL)
	//c.File("./build/index.html")
	c.JSON(http.StatusOK, gin.H{
		"error": "route not found",
	})
}

// DefineRoutes defines the routes for the gin engine.
// It sets up basic authentication and defines handlers for various routes.
// It also starts a goroutine to handle messages from the broadcast channel.
func DefineRoutes(appConfig config.AppConfig) *gin.Engine {
	go websocket2.WsManager.Start()
	appConf = appConfig

	router := gin.Default()
	router.RedirectTrailingSlash = false
	conf = appConfig

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "unity",
		"user":  "unity",
	}))
	router.GET(appConfig.BasePath+"/", handleRoot)
	router.GET(appConfig.BasePath+"/ping", handlePing)
	authorized.StaticFS(appConfig.BasePath+"/ui", http.Dir("./build"))
	authorized.GET(appConfig.BasePath+"/ws", handleWebsocket)
	router.GET(appConfig.BasePath+"/debug/pprof/*profile", gin.WrapF(pprof.Index))
	//router.Use(EnsureTrailingSlash())
	router.Use(LoggingMiddleware())
	router.Use(ErrorHandlingMiddleware())

	go func() {
		err := handleMessages()
		if err != nil {
			log.WithError(err).Error("Go routine crashed handling messages")
		}
	}()

	router.NoRoute(handleNoRoute)

	//processes.RunSync()
	return router
}

// handleMessages reads messages from the broadcast channel and handles them based on their type.
// It creates a new datastore and uses it to handle install messages.
func handleMessages() error {
	log.Info("Creating message handler")
	store, err := database.NewGormDatastore()
	if err != nil {
		log.WithError(err).Error("Error creating datastore")
		return err
	}

	for message := range websocket2.WsManager.Broadcast {
		// Unmarshal the message into a WebsocketMessage
		clientMessage := &marketplace.UnityWebsocketMessage{}
		if err := proto.Unmarshal(message.Message, clientMessage); err != nil {
			log.WithError(err).Error("Error unmarshalling websocket message")
			continue
		}

		log.Infof("Message recieved: %v", clientMessage)
		switch content := clientMessage.Content.(type) {
		case *marketplace.UnityWebsocketMessage_Install:
			installMessage := content.Install
			// Handle install message
			if err := processes.TriggerInstall(websocket2.WsManager, message.Client.UserID, store, installMessage, &conf); err != nil {
				log.WithError(err).Error("Error triggering install")
			}
		case *marketplace.UnityWebsocketMessage_Simplemessage:
			simpleMessage := content.Simplemessage
			resp, err := processes.ProcessSimpleMessage(simpleMessage, &conf, store, websocket2.WsManager, message.Client.UserID)
			if err != nil {
				log.WithError(err).Error("Problems parsing simple message")
			}
			websocket2.WsManager.SendMessageToClient(message.Client, resp)
		case *marketplace.UnityWebsocketMessage_Parameters:
			params := content.Parameters
			processes.UpdateParameters(params, store, &conf, websocket2.WsManager, message.Client.UserID)
		default:
			log.Error("Unknown message type")
		}
	}

	return nil
}
