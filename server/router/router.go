package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	chat "chat-app/internal/chat"
	"chat-app/internal/user"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, wsv1Handler *chat.Handler) {
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

	r.GET("/chat/v1", wsv1Handler.StartWS)
	r.GET("/chat/v1/get-clients", wsv1Handler.GetClients)
}

func Start(addr string) error {
	return r.Run(addr)
}
