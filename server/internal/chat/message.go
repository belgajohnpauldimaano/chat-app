package chat

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type MessageEvent struct {
	Type           string    `json:"type"`
	Sender         string    `json:"sender"`
	Recipient      string    `json:"recipient"`
	Content        string    `json:"content"`
	ContentType    int32     `json:"contentType"`
	ConversationID string    `json:"conversationId"`
	Timestamp      time.Time `json:"timestamp"`
}

type MessageEventHandler func(event *MessageEvent, h *Hub) error

const (
	PRIVATE_MESSAGE_EVENT           = "privateMessageEvent"
	PRIVATE_MESSAGE_EVENT_PUBLISHER = "privateMessageEventPublisher"
)

func SendPrivateMessagePublisher(messageEvent *MessageEvent, h *Hub) error {
	cachingClient := h.caching.RedisClientRing
	log.Println("Publishing a private message to pubsub...")

	jsonData, err := json.Marshal(messageEvent)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	payload := string(jsonData)

	// Publish a message to the channel
	ctx := context.Background()
	errPublish := cachingClient.Publish(ctx, PRIVATE_MESSAGE_EVENT, payload).Err()
	if errPublish != nil {
		log.Println("Error publishing message: ", errPublish)
		return errPublish
	}

	return nil
}

func SendPrivateMessage(messageEvent *MessageEvent, h *Hub) error {
	if client, ok := h.Clients[messageEvent.Recipient]; ok {
		log.Println("Sending private message to: ", messageEvent.Recipient)
		// single use can have multipe instances of connection when
		// logging in from mulple devices/browsers
		for _, c := range client {
			c.MessageEvent <- messageEvent
		}
	}
	return nil
}
