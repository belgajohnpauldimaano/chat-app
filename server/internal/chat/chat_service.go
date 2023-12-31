package chat

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

func (cs *chatService) GetConversations(ctx context.Context, request *GetUserConversationsRequest) ([]*Conversation, error) {
	return cs.chatRepository.GetConversations(ctx, request.UserId)
}

// TODO: Change the function name into Generate Message
// it includes:
//   - Create Conversation if needed
//   - Create message
func (cs *chatService) CreateConversation(ctx context.Context, conversation *ConversationRequest) (*Message, error) {
	log.Println("Service Creating conversation...")
	log.Println("at Service Converstion ID", conversation.ConversationId)

	existingUserConversation, existingUserConversationErr := cs.chatRepository.GetUsersExistingConversation(ctx, []string{"844b6164-3fb7-4197-90ba-625a719a2647", "aa44324d-bcd4-4e94-a0f7-ae0999e29bd1"})

	if existingUserConversationErr != nil {
		log.Println("Error in getting user conversation. err: ", existingUserConversationErr)
	}

	conversation.ConversationId = existingUserConversation.ConversationId
	log.Println("conversation.ConversationId: ", conversation.ConversationId, "xxxx")
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
		conversation.ConversationId = senderConverstionResult.ConversationId
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
	log.Println("conversation.ConversationId: ", conversation.ConversationId)
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

func (cs *chatService) GetMessagesByConversation(ctx context.Context, conversationId string) ([]*Message, error) {
	return cs.chatRepository.GetMessagesByConversation(ctx, conversationId)
}

func (cs *chatService) CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error) {
	return cs.chatRepository.CreateMessage(ctx, message)
}
