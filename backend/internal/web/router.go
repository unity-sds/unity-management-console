package web

import (
	"encoding/json"
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
	router.GET("/", func(c *gin.Context) {
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
		c.JSON(200, conf)
	})
	authorized.GET("/ws", func(c *gin.Context) {
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
				runner.UpdateCoreConfig(conn, store, conf)
			} else if received.Action == "install software" {
				//runner := &processes.ActRunnerImpl{}
				//if err != nil {
				//	log.Errorf("Failed to decode payload:", err)
				//	return
				//}
				//pb := &marketplace.Install{}
				//
				//err = proto.Unmarshal(received.Payload, pb)
				//if err != nil {
				//	log.Println("Error during message unmarshalling:", err)
				//	break
				//}
				//log.Infof("Message decoded successfully, %v", &pb)
				//err = runner.TriggerInstall(conn, store, *pb, conf)
				//if err != nil {
				//	log.Errorf("Error running workflow: %v", err)
				//}
			} else if received.Action == "request config" {
				msg, err := fetchConfig()
				if err != nil {
					log.Errorf("Problem requesting config: %v", err)
				}
				log.Info("Writing config to websocket")
				if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
					log.Errorf("Issue writing websocket message: %v", err)
					break
				}
			}

			// Echo the message back to the client.
			//if err := conn.WriteMessage(msgType, msg); err != nil {
			//	log.Println("Error during message writing:", err)
			//	break
			//}
		}
	})

	router.NoRoute(func(c *gin.Context) { // fallback
		c.File("./build/index.html")
	})

	return router
}

func fetchConfig() ([]byte, error) {

	pub, priv, err := aws.FetchSubnets()
	if err != nil {
		log.Errorf("error fetching subnets: %v", err)
	}
	netconfig := marketplace.Config_NetworkConfig{
		Publicsubnets:  pub,
		Privatesubnets: priv,
	}
	genconfig := &marketplace.Config{
		ApplicationConfig: nil,
		NetworkConfig:     &netconfig,
	}

	log.Infof("Config Generated: %+v", genconfig)
	return proto.Marshal(genconfig)
}

func decodeProtobuf(msg []byte, conn *websocket.Conn, conf config.AppConfig, store database.Datastore) error {
	pb := &marketplace.Install{}
	runner := &processes.ActRunnerImpl{}

	log.Infof("Decoding message: %v", msg)
	err := proto.Unmarshal(msg, pb)
	if err != nil {
		log.Println("Error during message unmarshalling:", err)
		return err
	}
	log.Infof("Message decoded successfully, %+v", pb)
	err = runner.TriggerInstall(conn, store, *pb, conf)
	if err != nil {
		log.Errorf("Error running workflow: %v", err)
	}
	return err
}
