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

// TODO: Change the function name into Generate Message
// it includes:
//   - Create Conversation if needed
//   - Create message
func (cs *chatService) CreateConversation(ctx context.Context, conversation *ConversationRequest) (*Message, error) {
	log.Println("Service Creating conversatin...")

	// Conversation is not created yet, so we create one
	if conversation.ConversationId == "" {
		senderConversation := &Conversation{
			ConversationType: conversation.ConversationType,
			UserId:           conversation.UserId,
		}

		recipientConversation := &Conversation{
			ConversationType: conversation.ConversationType,
			UserId:           conversation.RecipientId,
		}

		conversations := []*Conversation{}
		conversations = append(conversations, recipientConversation)

		senderConverstionResult, err := cs.chatRepository.CreateConversation(context.TODO(), senderConversation)
		log.Println("Conversation create: ", senderConverstionResult.ConversationId)
		if err != nil {
			log.Println("Error while creating conversation", err)
		}

		for _, conversation := range conversations {
			conversation.ConversationId = senderConverstionResult.ConversationId
			converstionResult, err := cs.chatRepository.CreateConversation(context.TODO(), conversation)
			log.Println("Conversation create: ", converstionResult.ConversationId)
			if err != nil {
				log.Println("Error while creating conversation", err)
			}
		}

		// return senderConversation, nil
	}

	messageRequest := &MessageRequest{
		ConversationId: conversation.ConversationId,
		SenderId:       conversation.UserId,
		RecipientId:    conversation.RecipientId,
		Content:        conversation.Content,
		ContentType:    conversation.ContentType,
	}

	newMessage, err := cs.chatRepository.CreateMessage(ctx, messageRequest)

	if err != nil {
		log.Println("Error while creating a message in database, err: ", err)
	}

	log.Println("Created time: ", newMessage.Timestamp)

	return newMessage, nil
}

func (cs *chatService) GetMessagesByConversation(ctx context.Context, conversationId int64) ([]*Message, error) {
	return cs.chatRepository.GetMessagesByConversation(ctx, conversationId)
}

func (cs *chatService) CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error) {
	return cs.chatRepository.CreateMessage(ctx, message)
}
