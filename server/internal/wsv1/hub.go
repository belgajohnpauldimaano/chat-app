package wsv1

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

	caching "chat-app/infrastructure/cache"
)

type Hub struct {
	Clients                      map[string][]*Client
	Register                     chan *Client
	Unregister                   chan *Client
	Broadcast                    chan *MessageEvent
	PrivateMessageEventBroadcast chan *MessageEvent
	caching                      *caching.RedisImpl
	messageEventHandlers         map[string]MessageEventHandler
	chatService                  ChatService

	sync.RWMutex
}

func NewHub(caching *caching.RedisImpl, chatService ChatService) *Hub {
	hub := &Hub{
		Clients:                      make(map[string][]*Client),
		Register:                     make(chan *Client),
		Unregister:                   make(chan *Client),
		Broadcast:                    make(chan *MessageEvent, 5),
		PrivateMessageEventBroadcast: make(chan *MessageEvent, 5),
		messageEventHandlers:         make(map[string]MessageEventHandler),
		caching:                      caching,
		chatService:                  chatService,
	}
	log.Println("NeHub")
	hub.setupEventHandlers()

	for i, x := range hub.messageEventHandlers {
		log.Println(i, x)
	}
	log.Println("NeHub")
	return hub
}

func (h *Hub) Run() {
	cachingClient := h.caching.RedisClientRing

	// Create a context with timeout
	ctx := context.Background()

	// Subscribe to a channel
	pubsub := cachingClient.Subscribe(ctx, PRIVATE_MESSAGE_EVENT)
	defer pubsub.Close()

	// Channel to handle messages
	messageChannel := pubsub.Channel()

	for {
		select {
		case pubsubMessage := <-messageChannel:
			// This is the redis pubsub handler for receiving data from redis
			// and publishing it to websocket
			log.Println("Received message:", pubsubMessage.Payload)

			var message MessageEvent

			if err := json.Unmarshal([]byte(pubsubMessage.Payload), &message); err != nil {
				log.Println("error unmarshalling request: ", err)
				continue
			}

			if err := h.routeEvent(&message, message.Type); err != nil {
				continue
			}

			// NOTE: This code below can send a data to MessageEvent Channel
			// if client, ok := h.Clients[message.Recipient]; ok {
			// 	log.Println("Sending private message to: ", message.Recipient)
			// 	// single use can have multipe instances of connection when
			// 	// logging in from mulple devices/browsers
			// 	for _, c := range client {
			// 		c.MessageEvent <- &message
			// 	}
			// }
		case client := <-h.Register:
			h.addClient(client)
		case client := <-h.Unregister:
			h.removeClient(client)
			// NOTE: THIS BELOW CAN BE REMOVED BECAUSE OF REDIS PUBSUB
			// USE THIS SETUP IF THERE IS NO INCLUDED REDIS
			// case message := <-h.Broadcast:
			// 	log.Println("Message type: ", message.Type)
			// 	log.Println("Message sender: ", message.Sender)
			// 	log.Println("Message recipient: ", message.Recipient)
			// 	log.Println("Message content: ", message.Content)
			// 	for _, client := range h.Clients {
			// 		for _, c := range client {
			// 			c.MessageEvent <- message
			// 		}
			// 	}
			// case message := <-h.PrivateMessageEventBroadcast:
			// 	log.Println("Message type: ", message.Type)
			// 	log.Println("Message sender: ", message.Sender)
			// 	log.Println("Message recipient: ", message.Recipient)
			// 	log.Println("Message content: ", message.Content)

			// 	jsonData, err := json.Marshal(message)
			// 	if err != nil {
			// 		log.Println("Error:", err)
			// 		return
			// 	}

			// 	log.Println(string(jsonData))
			// 	payload := string(jsonData)

			// 	// Publish a message to the channel
			// 	errPublish := cachingClient.Publish(ctx, PRIVATE_MESSAGE_EVENT, payload).Err()
			// 	if errPublish != nil {
			// 		log.Println("Error publishing message:", errPublish)
			// 		return
			// 	}

			// 	// TODO: Might need to remove this one because
			// 	// redis pubsub (pubsubMessage) will be
			// 	// added to the layer for resiliency
			// 	if client, ok := h.Clients[message.Recipient]; ok {
			// 		log.Println("Sending private message to: ", message.Recipient)
			// 		// single use can have multipe instances of connection when
			// 		// logging in from mulple devices/browsers
			// 		for _, c := range client {
			// 			c.MessageEvent <- message
			// 		}
			// 	}
			// NOTE: UP TO THIS CAN BE REMOVED
		}
	}
}

func (h *Hub) addClient(c *Client) {
	h.Lock()
	defer h.Unlock()

	log.Println("Adding Client to the client Pool...")

	if _, ok := h.Clients[c.UserId]; !ok {
		h.Clients[c.UserId] = make([]*Client, 0)
	}

	h.Clients[c.UserId] = append(h.Clients[c.UserId], c)
	log.Println(len(h.Clients[c.UserId]))

	log.Println("Client added to the client Pool...")
}

func (h *Hub) removeClient(c *Client) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.Clients[c.UserId]; ok {
		log.Println("Unregistering client...")
		delete(h.Clients, c.ID)
		close(c.MessageEvent)
		log.Println("Successfully unregistered client...")
	}
}

func (h *Hub) setupEventHandlers() {
	h.messageEventHandlers[PRIVATE_MESSAGE_EVENT] = SendPrivateMessage
	h.messageEventHandlers[PRIVATE_MESSAGE_EVENT_PUBLISHER] = SendPrivateMessagePublisher
}

func (h *Hub) routeEvent(messageEvent *MessageEvent, eventType string) error {
	eventHandler, ok := h.messageEventHandlers[eventType]
	if !ok {
		return errors.New("No handler for event type: " + messageEvent.Type)
	}

	if err := eventHandler(messageEvent, h); err != nil {
		log.Println("Error encountered: ", err)
		return err
	}

	return nil
}
