package router

import (
	// "server/internal/ws"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"chat-app/internal/user"
	"chat-app/router/ws"
	"chat-app/router/wsv2"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler, wsv2Handler *wsv2.Handler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	// health-check
	r.GET("/health-check", func(c *gin.Context) {
		healthCheckRes := map[string]interface{}{
			"status":  "OK",
			"version": "0.1.1",
		}
		c.JSON(http.StatusOK, healthCheckRes)
	})

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	r.POST("/ws/create-room", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/getClients/:roomId", wsHandler.GetClients)

	r.GET("/ws/v2", wsv2Handler.StartWS)
	r.GET("/ws/v2/get-clients", wsv2Handler.GetClients)
}

func Start(addr string) error {
	return r.Run(addr)
}
