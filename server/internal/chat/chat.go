package chat

import (
	"context"
	"time"
)

type Conversation struct {
	ID               int64     `json:"id" db:"id"`
	UserId           string    `json:"userId" db:"user_id"`
	ConversationId   string    `json:"conversationId" db:"conversation_id"`
	ConversationType int32     `json:"conversationType" db:"conversation_type"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
}

type ConversationRequest struct {
	ConversationId   string `json:"ConversationId"`
	UserId           string `json:"userId"`
	RecipientId      string `json:"recipientId"`
	Content          string `json:"content"`
	ContentType      int32  `json:"contentType"`
	ConversationType int32  `json:"conversationType"`
}

type Message struct {
	ID             string    `json:"id"`
	ConversationId string    `json:"conversationId"`
	SenderId       string    `json:"senderId"`
	RecipientId    string    `json:"recipientId"`
	Content        string    `json:"content"`
	ContentType    int32     `json:"contentType"`
	Timestamp      time.Time `json:"createdAt"`
}

type UserConversation struct {
	UserId         string    `json:"userId" db:"user_id"`
	ConversationId string    `json:"conversationId"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
}

type MessageRequest struct {
	ConversationId string `json:"conversationId"`
	SenderId       string `json:"senderId"`
	RecipientId    string `json:"recipientId"`
	Content        string `json:"content"`
	ContentType    int32  `json:"contentType"`
}

type ChatRepository interface {
	GetConversations(ctx context.Context, userId string) ([]*Conversation, error)
	CreateConversation(ctx context.Context, conversation *Conversation) (*Conversation, error)
	GetMessagesByConversation(ctx context.Context, conversationId string) ([]*Message, error)
	CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error)
	GetUsersExistingConversation(ctx context.Context, userIds []string) (*Conversation, error)
}

type ChatService interface {
	GetConversations(ctx context.Context, request *GetUserConversationsRequest) ([]*Conversation, error)
	CreateConversation(ctx context.Context, conversations *ConversationRequest) (*Message, error)
	GetMessagesByConversation(ctx context.Context, conversationId string) ([]*Message, error)
	CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error)
}

type GetUserConversationsRequest struct {
	UserId string `json:"userId"`
	Page   int32  `json:"page"`
	Limit  int32  `json:"limit"`
}

type GetUserConversationsResponse struct {
	IsSuccess     bool            `json:"isSuccess"`
	Status        int32           `json:"status"`
	Message       string          `json:"message"`
	Conversations []*Conversation `json:"conversations"`
}
