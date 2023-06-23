package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/aws"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
	"github.com/unity-sds/unity-control-plane/backend/internal/processes"
	ws "github.com/unity-sds/unity-control-plane/backend/internal/websocket"
	"net/http"
)

var conf config.AppConfig
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // or check the origin if you want to add more security
	},
}

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

func handleRoot(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/ui")
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

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
	runner := &processes.ActRunnerImpl{}
	if err := processes.UpdateCoreConfig(nil, store, conf, runner); err != nil {
		log.WithError(err).Error("error updating core configuration")
	}
}

func handleConfigGET(c *gin.Context) {
	c.JSON(http.StatusOK, conf)
}

func handleWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	store := database.GormDatastore{}
	if err != nil {
		log.Print("upgrade error:", err)
		return
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}

		var received ws.BareMessage
		err = json.Unmarshal(msg, &received)
		if err != nil {
			log.Println("Error during message unmarshalling:", err)
			err := decodeProtobuf(msg, conn, conf, store)
			if err != nil {
				log.Errorf("Error decoding protobuf %v", err)
			}
			break
		}

		log.Infof("Message received : %v", received.Payload)
		log.Infof("Action received: %v", received.Action)
		if received.Action == "config upgrade" {
			runner := &processes.ActRunnerImpl{}
			processes.UpdateCoreConfig(conn, store, conf, runner)
		} else if received.Action == "install software" {

		} else if received.Action == "request config" {
			msg, err := fetchConfig(conf)
			if err != nil {
				log.Errorf("Problem requesting config: %v", err)
			}
			log.Info("Writing config to websocket")
			if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				log.Errorf("Issue writing websocket message: %v", err)
				break
			}
		} else if received.Action == "request parameters" {
			existingparams, err := store.FetchSSMParams()

			params, err := aws.ReadSSMParameters(existingparams)

			if err != nil {
				log.Errorf("Problem requesting config: %v", err)
			}
			log.Info("Writing params to websocket")
			msg, err := proto.Marshal(params)
			if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				log.Errorf("Issue writing websocket message: %v", err)
				break
			}
		}

	}
}

func handleNoRoute(c *gin.Context) {
	c.File("./build/index.html")
}

func DefineRoutes(appConfig config.AppConfig) *gin.Engine {
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

	return router
}
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
		GithubToken: conf.GithubToken,
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

func decodeProtobuf(msg []byte, conn *websocket.Conn, conf config.AppConfig, store database.Datastore) error {
	pb := &marketplace.Install{}
	runner := processes.NewActRunner() // using a constructor function

	log.Info("Attempting to decode the message...")
	if err := proto.Unmarshal(msg, pb); err != nil {
		return fmt.Errorf("failed to unmarshal message: %v, error: %w", msg, err)
	}

	log.Infof("Message decoded successfully: %+v", pb)

	if err := processes.TriggerInstall(conn, store, *pb, conf, *runner); err != nil {
		return fmt.Errorf("failed to trigger install: %v", err)
	}

	return nil
}
