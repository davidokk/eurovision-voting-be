-- Remove email-based auth (replaced by Telegram)

DROP TABLE IF EXISTS auth_codes;

DROP INDEX IF EXISTS users_email_lower_idx;

ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;
ALTER TABLE users DROP COLUMN IF EXISTS email;
