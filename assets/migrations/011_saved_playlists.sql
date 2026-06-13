CREATE TABLE IF NOT EXISTS saved_playlists (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name text NOT NULL,
    entries jsonb NOT NULL DEFAULT '[]'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS saved_playlists_user_idx ON saved_playlists (user_id, updated_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS saved_playlists_user_name_idx ON saved_playlists (user_id, lower(trim(name)));
