package wsv1

import (
	"context"
	"encoding/json"
	"log"
)

type MessageEvent struct {
	Type           string `json:"type"`
	Sender         string `json:"sender"`
	Recipient      string `json:"recipient"`
	Content        string `json:"content"`
	ConversationID string `json:"conversation_id"`
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

	log.Println(string(jsonData))
	payload := string(jsonData)
	log.Println("publishing payload: ", payload)

	// Create a conversation based on conversation type
	newConversation := &ConversationRequest{
		ID:               messageEvent.ConversationID,
		UserId:           messageEvent.Sender,
		RecipientId:      messageEvent.Recipient,
		ConversationType: 0,
	}
	h.chatService.CreateConversation(context.TODO(), newConversation)
	// Publish a message to the channel
	ctx := context.Background()
	errPublish := cachingClient.Publish(ctx, PRIVATE_MESSAGE_EVENT, payload).Err()
	if errPublish != nil {
		log.Println("Error publishing message:", errPublish)
		return errPublish
	}
	return nil
}

func SendPrivateMessage(messageEvent *MessageEvent, h *Hub) error {
	// TODO: Write a message to websocket here

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
