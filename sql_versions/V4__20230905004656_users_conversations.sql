CREATE TABLE users_conversations (
	user_id uuid NOT NULL,
	conversation_id uuid NOT NULL,
	created_at timestamp with time zone NULL DEFAULT now(),
	CONSTRAINT users_conversations_pk PRIMARY KEY (user_id,conversation_id)
);
