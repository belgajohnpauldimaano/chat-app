package chat

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type chatRepository struct {
	db DBTX
}

func NewChatRepository(db DBTX) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) GetConversations(ctx context.Context, userId int64) ([]*Conversation, error) {
	conversations := make([]*Conversation, 0)

	return conversations, nil
}

func (r *chatRepository) CreateConversation(ctx context.Context, conversation *Conversation) (*Conversation, error) {
	log.Println("Creating conversation...")
	log.Println("Sender: ", conversation.UserId)

	if conversation.ConversationId == "" {
		conversation.ConversationId = uuid.New().String()
	}

	result, err := r.db.ExecContext(ctx, `INSERT INTO conversations(conversation_id, conversation_type, user_id) VALUES($1, $2, $3)`, conversation.ConversationId, 0, conversation.UserId)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Unable to create conversation")
	}

	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", rows)
		return nil, errors.New("Unable to create conversation")
	}

	return conversation, nil
}

func (r *chatRepository) GetMessagesByConversation(ctx context.Context, conversationId int64) ([]*Message, error) {
	messages := make([]*Message, 0)

	return messages, nil
}

func (r *chatRepository) CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error) {
	log.Println("Creating message...")
	log.Println("Sender: ", message.SenderId)
	query := `
		INSERT INTO
			public.messages (id, conversation_id, sender_id, recipient_id, "content", content_type)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING "timestamp";`

	stmt, err := r.db.PrepareContext(ctx, query)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	genearatedMessageId := uuid.New().String()
	var createdTimestamp time.Time
	queryErr := stmt.QueryRow(genearatedMessageId, message.ConversationId, message.SenderId, message.RecipientId, message.Content, message.ContentType).Scan(&createdTimestamp)

	if queryErr != nil {
		log.Fatal(err)
		return nil, queryErr
	}

	createdMessage := &Message{
		ID:             genearatedMessageId,
		ConversationId: message.ConversationId,
		SenderId:       message.SenderId,
		RecipientId:    message.RecipientId,
		Content:        message.Content,
		ContentType:    message.ContentType,
		Timestamp:      createdTimestamp,
	}

	return createdMessage, nil
}