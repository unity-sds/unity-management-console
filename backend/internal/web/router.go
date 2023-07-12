package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"github.com/unity-sds/unity-control-plane/backend/internal/action"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/aws"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
	"github.com/unity-sds/unity-control-plane/backend/internal/processes"
	websocket2 "github.com/unity-sds/unity-control-plane/backend/internal/websocket"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"net/http"
)

var conf config.AppConfig

var wsManager = websocket2.NewWebSocketManager()

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan websocket2.ClientMessage)

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
	c.Redirect(http.StatusMovedPermanently, "/ui")
}

// handlePing responds with a JSON message containing "pong".
// This can be used to check if the server is running.
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// handleConfigPOST handles POST requests to "/config".
// It binds the JSON body of the request to a slice of CoreConfig models.
// If the binding is successful, it stores the configuration in the database and triggers an environment update.
func handleConfigPOST(c *gin.Context) {
	var configjson []models.CoreConfig
	store := database.GormDatastore{}

	if err := c.ShouldBindJSON(&configjson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := database.StoreConfig(configjson); err != nil {
		log.WithError(err).Error("error storing configuration")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": configjson})

	// Trigger environment update via act
	runner := &action.ActRunnerImpl{}
	if err := processes.UpdateCoreConfig(nil, store, conf, runner); err != nil {
		log.WithError(err).Error("error updating core configuration")
	}
}

// handleConfigGET responds with the current application configuration.
func handleConfigGET(c *gin.Context) {
	c.JSON(http.StatusOK, conf)
}

// handleWebsocket handles websocket connections.
// It upgrades the HTTP connection to a websocket connection and reads messages from the client.
// Each message is unmarshalled into a WebsocketMessage and sent to the broadcast channel.
func handleWebsocket(c *gin.Context) {
	wsManager.HandleConnections(c.Writer, c.Request)
}

// handleNoRoute serves the index.html file for any routes that are not defined.
func handleNoRoute(c *gin.Context) {
	c.File("./build/index.html")
}

// DefineRoutes defines the routes for the gin engine.
// It sets up basic authentication and defines handlers for various routes.
// It also starts a goroutine to handle messages from the broadcast channel.
func DefineRoutes(appConfig config.AppConfig) *gin.Engine {
	go wsManager.Start()

	router := gin.Default()
	conf = appConfig

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "unity",
		"user":  "unity",
	}))

	router.GET("/", handleRoot)
	router.GET("/ping", handlePing)
	authorized.StaticFS("/ui", http.Dir("./build"))
	authorized.POST("/config", handleConfigPOST)
	authorized.GET("/config", handleConfigGET)
	authorized.GET("/ws", handleWebsocket)
	router.NoRoute(handleNoRoute)

	router.Use(LoggingMiddleware())
	router.Use(ErrorHandlingMiddleware())

	go handleMessages()

	return router
}

// fetchConfig fetches the application and network configuration from AWS and marshals it into a protobuf message.
func fetchConfig(conf config.AppConfig) ([]byte, error) {

	pub, priv, err := aws.FetchSubnets()
	if err != nil {
		log.WithError(err).Error("Error fetching subnets")
		return nil, err
	}

	netconfig := marketplace.Config_NetworkConfig{
		Publicsubnets:  pub,
		Privatesubnets: priv,
	}

	appConfig := marketplace.Config_ApplicationConfig{
		GithubToken:      conf.GithubToken,
		MarketplaceOwner: conf.MarketplaceOwner,
		MarketplaceUser:  conf.MarketplaceRepo,
	}

	genconfig := &marketplace.Config{
		ApplicationConfig: &appConfig,
		NetworkConfig:     &netconfig,
	}

	log.WithFields(log.Fields{
		"Config": genconfig,
	}).Info("Config Generated")

	data, err := proto.Marshal(genconfig)
	if err != nil {
		log.WithError(err).Error("Failed to marshal config")
		return nil, err
	}

	return data, nil
}

// handleMessages reads messages from the broadcast channel and handles them based on their type.
// It creates a new datastore and uses it to handle install messages.
func handleMessages() error {
	store, err := database.NewGormDatastore()
	if err != nil {
		log.WithError(err).Error("Error creating datastore")
		return err
	}
	for message := range wsManager.Broadcast {
		// Unmarshal the message into a WebsocketMessage
		clientMessage := &marketplace.WebsocketMessage{}
		if err := proto.Unmarshal(message.Message, clientMessage); err != nil {
			log.WithError(err).Error("Error unmarshalling websocket message")
			continue
		}

		switch content := clientMessage.Content.(type) {
		case *marketplace.WebsocketMessage_Install:
			installMessage := content.Install
			// Handle install message
			if err := processes.TriggerInstall(wsManager, message.Client.UserID, store, installMessage, conf); err != nil {
				log.WithError(err).Error("Error triggering install")
			}
		default:
			log.Error("Unknown message type")
		}
	}

	return nil
}
