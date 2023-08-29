ALTER TABLE public.conversations ADD user_id uuid NULL;
ALTER TABLE public.conversations ADD CONSTRAINT conversations_user_fk FOREIGN KEY (user_id) REFERENCES public.users(id);
