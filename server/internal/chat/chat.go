package chat

import (
	"context"
	"time"
)

type Conversation struct {
	ID               int64     `json:"id"`
	ConversationId   string    `json:"ConversationId"`
	ConversationType int32     `json:"conversationType"`
	UserId           string    `json:"userId"`
	CreatedAt        time.Time `json:"createdAt"`
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

type MessageRequest struct {
	ConversationId string `json:"conversationId"`
	SenderId       string `json:"senderId"`
	RecipientId    string `json:"recipientId"`
	Content        string `json:"content"`
	ContentType    int32  `json:"contentType"`
}

type ChatRepository interface {
	GetConversations(ctx context.Context, userId int64) ([]*Conversation, error)
	CreateConversation(ctx context.Context, conversation *Conversation) (*Conversation, error)
	GetMessagesByConversation(ctx context.Context, conversationId int64) ([]*Message, error)
	CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error)
}

type ChatService interface {
	GetConversations(ctx context.Context, userId int64) ([]*Conversation, error)
	CreateConversation(ctx context.Context, conversations *ConversationRequest) (*Message, error)
	GetMessagesByConversation(ctx context.Context, conversationId int64) ([]*Message, error)
	CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error)
}
