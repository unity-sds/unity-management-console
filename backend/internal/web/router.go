package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"github.com/unity-sds/unity-control-plane/backend/internal/act"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
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
	ffclient := config.InitApplication()

	username := c.MustGet(gin.AuthUserKey).(string)
	user := ffuser.NewUser(username)

	hasFlag, _ := ffclient.BoolVariation("test-flag", user, false)
	if hasFlag { // flag "test-flag" is true for the user
		log.Info("Flag true")
	} else { // flag "test-flag" is false for the user
		log.Info("flag false")
	}
}
func DefineRoutes() *gin.Engine {
	router := gin.Default()
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "unity",
		"user": "unity",
	}))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authorized.StaticFS("/ui", http.Dir("./build"))

	authorized.POST("/config", func(c *gin.Context){
		var config models.CoreConfig

		if err := c.ShouldBindJSON(&config); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := database.DB.Create(&config).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": config})

	})

	authorized.GET("/config", func(c *gin.Context){

	})
	authorized.GET("/ws", func(c *gin.Context) {
		setupFeatureFlags(c)
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("upgrade error:", err)
			return
		}
		act.RunAct(conn)
	})

	router.NoRoute(func(c *gin.Context) { // fallback
		c.File("./build/index.html")
	})

	return router
}
