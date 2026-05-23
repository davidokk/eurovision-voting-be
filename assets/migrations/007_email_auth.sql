-- Email-based auth: unique login by email, verification codes

ALTER TABLE users ADD COLUMN IF NOT EXISTS email text;
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified_at timestamptz;

-- Existing users: placeholder until they bind a real email
UPDATE users
SET email = id::text || '@legacy.pending',
    email_verified_at = NULL
WHERE email IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_idx ON users (lower(trim(email)))
WHERE email IS NOT NULL;

CREATE TABLE IF NOT EXISTS auth_codes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email text NOT NULL,
    code_hash text NOT NULL,
    purpose text NOT NULL,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    username text,
    password_hash text,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    used_at timestamptz
);

CREATE INDEX IF NOT EXISTS auth_codes_active_idx
    ON auth_codes (lower(trim(email)), purpose)
    WHERE used_at IS NULL;
