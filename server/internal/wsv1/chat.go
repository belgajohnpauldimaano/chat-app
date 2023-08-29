package wsv1

import (
	"context"
	"time"
)

type Conversation struct {
	ID               string    `json:"id"`
	ConversationType int32     `json:"conversation_type"`
	UserId           string    `json:"user_id"`
	CreatedAt        time.Time `json:"created_at"`
}

type ConversationRequest struct {
	ID               string `json:"id"`
	UserId           string `json:"user_id"`
	RecipientId      string `json:"recipient_id"`
	ConversationType int32  `json:"conversation_type"`
}

type Message struct {
	ID             string    `json:"id"`
	ConversationId string    `json:"conversation_id"`
	SenderId       string    `json:"sender_id"`
	RecipientId    string    `json:"recipient_id"`
	Content        string    `json:"content"`
	ContentType    int64     `json:"content_type"`
	Timestamp      time.Time `json:"created_at"`
}

type MessageRequest struct {
	ConversationId int64  `json:"conversation_id"`
	SenderId       int64  `json:"sender_id"`
	RecipientId    int64  `json:"recipient_id"`
	Content        string `json:"content"`
	ContentType    int64  `json:"content_type"`
}

type ChatRepository interface {
	GetConversations(ctx context.Context, userId int64) ([]*Conversation, error)
	CreateConversation(ctx context.Context, conversation []*Conversation) (*Conversation, error)
	GetMessagesByConversation(ctx context.Context, conversationId int64) ([]*Message, error)
	CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error)
}

type ChatService interface {
	GetConversations(ctx context.Context, userId int64) ([]*Conversation, error)
	CreateConversation(ctx context.Context, conversations *ConversationRequest) (*Conversation, error)
	GetMessagesByConversation(ctx context.Context, conversationId int64) ([]*Message, error)
	CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error)
}
