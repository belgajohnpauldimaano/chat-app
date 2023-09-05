package chat

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub         *Hub
	chatService ChatService
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) StartWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.Query("userId")
	client := &Client{
		Conn:         conn,
		MessageEvent: make(chan *MessageEvent, 10),
		ID:           userId,
		UserId:       userId,
	}

	h.hub.Register <- client

	go client.writeMessage()
	client.readMessage(h.hub)
}

type ClientRes struct {
	ID     string `json:"id"`
	UserId string `json:"userId"`
}

func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientRes

	for _, c := range h.hub.Clients {
		cl := ClientRes{
			ID:     c[0].ID,
			UserId: c[0].UserId,
		}
		clients = append(clients, cl)
	}

	c.JSON(http.StatusOK, clients)
}

func (h *Handler) GetUserConversations(c *gin.Context) {
	var requestObj GetUserConversationsRequest
	if c.Query("page") != "" {
		page, err := strconv.ParseInt(c.Query("page"), 10, 32)
		if err == nil {
			requestObj.Page = int32(page)
		}
	}

	if c.Query("limit") != "" {
		limit, err := strconv.ParseInt(c.Query("limit"), 10, 32)
		if err == nil {
			requestObj.Limit = int32(limit)
		}
	}
	requestObj.UserId = c.Query("userId")
	conversations, getConversationsErr := h.hub.chatService.GetConversations(c.Request.Context(), &requestObj)

	responseObj := &GetUserConversationsResponse{
		IsSuccess:     true,
		Status:        200,
		Message:       "User conversations",
		Conversations: conversations,
	}

	if getConversationsErr != nil {
		responseObj.IsSuccess = false
		responseObj.Status = 500
		responseObj.Message = "Something went wrong!"
		responseObj.Conversations = nil
	}

	c.JSON(http.StatusOK, responseObj)
}
