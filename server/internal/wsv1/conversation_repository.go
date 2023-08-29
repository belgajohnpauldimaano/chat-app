package wsv1

import (
	"context"
	"database/sql"
	"log"

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

func (r *chatRepository) CreateConversation(ctx context.Context, conversations []*Conversation) (*Conversation, error) {
	log.Println("Creating conversatin...")
	log.Println("Sender: ", conversations[0].UserId)

	senderResult, err := r.db.ExecContext(ctx, `INSERT INTO conversations(id, conversation_type, user_id) VALUES($1, $2, $3)`, uuid.New().String(), 0, conversations[0].UserId)
	if err != nil {
		log.Fatal(err)
	}
	senderRows, err := senderResult.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if senderRows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", senderRows)
	}

	recipientResult, err := r.db.ExecContext(ctx, `INSERT INTO conversations(id, conversation_type, user_id) VALUES($1, $2, $3)`, uuid.New().String(), 0, conversations[0].UserId)
	if err != nil {
		log.Fatal(err)
	}

	recipientRows, err := recipientResult.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if recipientRows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", recipientRows)
	}

	createdConversation := &Conversation{}

	return createdConversation, nil
}

func (r *chatRepository) GetMessagesByConversation(ctx context.Context, conversationId int64) ([]*Message, error) {
	messages := make([]*Message, 0)

	return messages, nil
}

func (r *chatRepository) CreateMessage(ctx context.Context, message *MessageRequest) (*Message, error) {
	createdMessage := &Message{}

	return createdMessage, nil
}
