package chat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
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

func (r *chatRepository) GetConversations(ctx context.Context, userId string) ([]*Conversation, error) {
	var conversations []*Conversation
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT * FROM conversations
			WHERE user_id = $1
			LIMIT 20
		`,
		userId,
	)

	if err != nil {
		log.Println("Something went wrong, error: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var result Conversation
		err := rows.Scan(&result.ID, &result.ConversationId, &result.UserId, &result.ConversationType, &result.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		conversations = append(conversations, &result)
	}

	return conversations, nil
}

func (r *chatRepository) CreateConversation(ctx context.Context, conversation *Conversation) (*Conversation, error) {
	log.Println("Creating conversation...")
	log.Println("Sender: ", conversation.UserId)
	log.Println("convo id: ", conversation.ConversationId)
	log.Println("Type: ", conversation.ConversationType)

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

func (r *chatRepository) GetMessagesByConversation(ctx context.Context, conversationId string) ([]*Message, error) {
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
		log.Println("Something went wrong! err: ", queryErr)
		log.Fatal(queryErr)
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

func (r *chatRepository) GetUsersExistingConversation(ctx context.Context, userIds []string) (*Conversation, error) {
	var inParameterValue string = ""

	for _, userId := range userIds {
		inParameterValue = fmt.Sprintf(`%s'%s',`, inParameterValue, userId)
	}

	inParameterValue = inParameterValue[:len(inParameterValue)-1]

	usersConversationsQuery := fmt.Sprintf(`
		select conversation_id from conversation_participants
		where user_id in (%s)
		group by conversation_id 
		having  count(distinct user_id) = %s;
		`, inParameterValue, strconv.Itoa(len(userIds)))

	usersConversationsRows, err := r.db.QueryContext(
		ctx,
		usersConversationsQuery,
	)

	defer usersConversationsRows.Close()

	if err != nil {
		log.Println("Something went wrong, error: ", err)
		return nil, err
	}

	var userConversation UserConversation

	for usersConversationsRows.Next() {
		err := usersConversationsRows.Scan(&userConversation.ConversationId)
		if err != nil {
			log.Println("Scanning error: ", err)
		}
	}

	log.Println("UserId: ", userConversation.UserId)
	log.Println("ConversationId: ", userConversation.ConversationId)

	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT * FROM conversations
			WHERE conversation_id = $1
			LIMIT 20
		`,
		userConversation.ConversationId,
	)

	if err != nil {
		log.Println("Something went wrong, error: ", err)
		return nil, err
	}

	defer rows.Close()

	var result Conversation
	rows.Next()
	err = rows.Scan(&result.ID, &result.UserId, &result.ConversationId, &result.ConversationType, &result.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}

	return &result, nil
}
