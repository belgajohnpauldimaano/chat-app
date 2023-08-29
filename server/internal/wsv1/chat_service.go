package wsv1

import (
	"context"
	"log"
)

type chatService struct {
	chatRepository ChatRepository
}

func NewChatService(chatRepository ChatRepository) ChatService {
	return &chatService{
		chatRepository: chatRepository,
	}
}

func (cs *chatService) GetConversations(ctx context.Context, userId int64) ([]*Conversation, error) {
	return cs.chatRepository.GetConversations(ctx, userId)
}

func (cs *chatService) CreateConversation(ctx context.Context, conversation *ConversationRequest) (*Conversation, error) {
	log.Println("Service Creating conversatin...")

	senderConversation := &Conversation{
		ConversationType: conversation.ConversationType,
		UserId:           conversation.UserId,
	}

	recipientConversation := &Conversation{
		ConversationType: conversation.ConversationType,
		UserId:           conversation.RecipientId,
	}

	conversations := []*Conversation{}
	conversations = append(conversations, senderConversation, recipientConversation)

	return cs.chatRepository.CreateConversation(context.TODO(), conversations)
}

func (cs *chatService) GetMessagesByConversation(ctx context.Context, conversationId int64) ([]*Message, error) {
	return cs.chatRepository.GetMessagesByConversation(ctx, conversationId)
}

func (cs *chatService) CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error) {
	return cs.chatRepository.CreateMessage(ctx, message)
}
