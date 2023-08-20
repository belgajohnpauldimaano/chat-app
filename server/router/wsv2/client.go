package wsv2

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID           string `json:"id"`
	UserId       string `json:"userId"`
	Conn         *websocket.Conn
	MessageEvent chan *MessageEvent
}

// type MessageContent struct {

// }

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10 // 90% of pongWait
)

const (
	privateMessage = "privateMessage"
	groupMessage   = "groupMessage"
)

// type Payload struct {
// 	Recipient string `json:"recipient"`
// 	Content   string `json:"content"`
// }

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("Error setting read deadline: ", err)
		return
	}

	c.Conn.SetPongHandler(c.PongHandler)
	c.Conn.SetReadLimit(512)

	for {

		messageType, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		log.Println("Message type: ", messageType)
		log.Println("Payload: ", payload)

		var messageEventContent MessageEvent

		if err := json.Unmarshal(payload, &messageEventContent); err != nil {
			log.Println("error unmarshalling request: ", err)
			continue
		}
		messageEventContent.Sender = c.UserId

		// msg := &MessageEvent{
		// 	Type:      "test",
		// 	Sender:    c.UserId,
		// 	Recipient: "",
		// 	Content:   string(payload),
		// }

		// msg := &messageEventContent

		switch messageEventContent.Type {
		case privateMessage:
			log.Println("private message")
			hub.PrivateMessageEventBroadcast <- &messageEventContent
		default:
			hub.MessageEventBroadcast <- &messageEventContent
		}
	}
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.MessageEvent:
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection close: ", err)
					return
				}
			}

			data, err := json.Marshal(message)

			if err != nil {
				log.Println("Erro marshalling message: ", err)
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("Failed to send message: ", err)
			}

			log.Println("Message sent: ", message)
			continue
		case <-ticker.C:
			log.Println("Sending ping...")
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Failed to send ping: ", err)
				return
			}
		}
	}
}

func (c *Client) PongHandler(pongMsg string) error {
	log.Println("PongHandler: ", pongMsg)
	return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
}
