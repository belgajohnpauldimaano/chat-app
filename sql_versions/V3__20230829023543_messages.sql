CREATE TABLE "messages" (
	"id" uuid PRIMARY KEY,
  "conversation_id" uuid NOT NULL,
  "sender_id" uuid NOT NULL,
  "recipient_id" uuid NOT NULL,
  "content" text,
  "content_type" int NOT NULL,
  "timestamp" TIMESTAMP DEFAULT NOW(),
	CONSTRAINT message_sender_user_user_id_fk FOREIGN KEY (sender_id) REFERENCES users(id),
	CONSTRAINT messages_recipient_user_user_id_fk FOREIGN KEY (recipient_id) REFERENCES users(id)
)
