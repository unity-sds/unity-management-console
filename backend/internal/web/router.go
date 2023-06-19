package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
	"github.com/unity-sds/unity-control-plane/backend/internal/processes"
	"net/http"
)

type Message struct {
	Action  string              `json:"action"`
	Payload []models.CoreConfig `json:"payload"`
}

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
func DefineRoutes(conf config.AppConfig) *gin.Engine {
	router := gin.Default()
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "unity",
		"user":  "unity",
	}))
	router.GET("/", func(c *gin.Context){
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authorized.StaticFS("/ui", http.Dir("./build"))

	authorized.POST("/config", func(c *gin.Context) {
		// Persist settings
		var configjson []models.CoreConfig

		store := database.GormDatastore{}
		if err := c.ShouldBindJSON(&configjson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := database.StoreConfig(configjson)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": configjson})

		// Trigger environment update via act
		runner := &processes.ActRunnerImpl{}
		err = runner.UpdateCoreConfig(nil, store, conf)
		if err != nil {
			return
		}
	})

	authorized.GET("/config", func(c *gin.Context) {

	})
	authorized.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		store := database.GormDatastore{}
		if err != nil {
			log.Print("upgrade error:", err)
			return
		}
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error during message reading:", err)
				break
			}

			var received Message
			err = json.Unmarshal(msg, &received)
			if err != nil {
				log.Println("Error during message unmarshalling:", err)
				break
			}

			log.Infof("Message received : %v", received.Payload)
			log.Infof("Action received: %v", received.Action)
			if received.Action == "config upgrade" {
				runner := &processes.ActRunnerImpl{}
				runner.UpdateCoreConfig(conn, store, conf)
			} else if received.Action == "install software" {
				runner := &processes.ActRunnerImpl{}
				err := runner.InstallMarketplaceApplication(conn, store, received.Payload[0].Value, conf)
				if err != nil{
					log.Errorf("Error running workflow: %v", err)
				}
			}

			// Echo the message back to the client.
			if err := conn.WriteMessage(msgType, msg); err != nil {
				log.Println("Error during message writing:", err)
				break
			}
		}
	})

	router.NoRoute(func(c *gin.Context) { // fallback
		c.File("./build/index.html")
	})

	return router
}
