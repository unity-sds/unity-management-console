package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/act"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // or check the origin if you want to add more security
	},
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

	authorized.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("upgrade error:", err)
			return
		}
		act.RunAct(conn)
	})
	//router.Use(static.Serve("/", static.LocalFile("./build", true)))


	router.NoRoute(func(c *gin.Context) { // fallback
		c.File("./build/index.html")
	})

	return router
}
