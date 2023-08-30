CREATE TABLE "conversations" (
	"id" bigserial PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "conversation_id" uuid NOT NULL,
  "conversation_type" int NOT NULL, --- This column is to distinguish private chat as 0 or group chat as 1
  "created_at" TIMESTAMP DEFAULT NOW(),
  CONSTRAINT conversation_user_fk FOREIGN KEY (user_id) REFERENCES users(id)
);
