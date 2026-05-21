-- Chat media, replies, avatars
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url text;

ALTER TABLE messages ADD COLUMN IF NOT EXISTS id uuid;
UPDATE messages SET id = gen_random_uuid() WHERE id IS NULL;
ALTER TABLE messages ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE messages ALTER COLUMN id SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'messages_pkey'
  ) THEN
    ALTER TABLE messages ADD PRIMARY KEY (id);
  END IF;
END $$;

ALTER TABLE messages ADD COLUMN IF NOT EXISTS reply_to_id uuid REFERENCES messages(id) ON DELETE SET NULL;
ALTER TABLE messages ADD COLUMN IF NOT EXISTS content_type text NOT NULL DEFAULT 'text';
ALTER TABLE messages ADD COLUMN IF NOT EXISTS media_url text;
ALTER TABLE messages ADD COLUMN IF NOT EXISTS media_duration_ms integer;

CREATE INDEX IF NOT EXISTS messages_contest_created_idx ON messages (contest_id, created_at);
CREATE INDEX IF NOT EXISTS messages_reply_to_idx ON messages (reply_to_id);
