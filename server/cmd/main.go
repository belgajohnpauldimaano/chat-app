package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/cache/v9"

	caching "chat-app/infrastructure/cache"
	"chat-app/infrastructure/db"
	chat "chat-app/internal/chat"
	"chat-app/internal/user"
	"chat-app/router"
)

func main() {

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	defer dbConn.Close()
	// start

	type Object struct {
		Str string
		Num int
	}

	redisClient := caching.NewRedisClient()
	// redisClient.Cache().Get()
	defer redisClient.Close()

	ctx := context.TODO()
	key := "mykey"
	obj := &Object{
		Str: "mystring",
		Num: 42,
	}

	if err := redisClient.CacheClient.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		panic(err)
	}

	var wanted Object
	if err := redisClient.CacheClient.Get(ctx, key, &wanted); err == nil {
		fmt.Println(wanted)
	}
	// End

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	chatRepo := chat.NewChatRepository(dbConn.GetDB())
	chatService := chat.NewChatService(chatRepo)
	hubV1 := chat.NewHub(redisClient, chatService)
	wsHandlerV1 := chat.NewHandler(hubV1)

	go hubV1.Run()

	router.InitRouter(userHandler, wsHandlerV1)
	// For Redis Pubsub testing on multiple app instance
	// First instance
	router.Start("0.0.0.0:8080")
	// Second instance
	// router.Start("0.0.0.0:8081")
}
