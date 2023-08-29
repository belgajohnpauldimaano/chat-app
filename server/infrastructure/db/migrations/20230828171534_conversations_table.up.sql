CREATE TABLE "conversations" (
	"id" uuid PRIMARY KEY,
  "conversation_type" int NOT NULL, --- This column is to distinguish private chat as 0 or group chat as 1
  "created_at" TIMESTAMP DEFAULT NOW()
);
