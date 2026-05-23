-- Telegram account linking and auth sessions

ALTER TABLE users ADD COLUMN IF NOT EXISTS telegram_id bigint;
ALTER TABLE users ADD COLUMN IF NOT EXISTS telegram_username text;
ALTER TABLE users ADD COLUMN IF NOT EXISTS telegram_linked_at timestamptz;

CREATE UNIQUE INDEX IF NOT EXISTS users_telegram_id_idx ON users (telegram_id) WHERE telegram_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS telegram_auth_sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    link_token text NOT NULL UNIQUE,
    purpose text NOT NULL,
    username text,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    telegram_id bigint,
    telegram_chat_id bigint,
    telegram_username text,
    code_hash text,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    used_at timestamptz,
    code_sent_at timestamptz
);

CREATE INDEX IF NOT EXISTS telegram_auth_sessions_active_token_idx
    ON telegram_auth_sessions (link_token)
    WHERE used_at IS NULL;

CREATE INDEX IF NOT EXISTS telegram_auth_sessions_telegram_rate_idx
    ON telegram_auth_sessions (telegram_id, created_at)
    WHERE telegram_id IS NOT NULL;
