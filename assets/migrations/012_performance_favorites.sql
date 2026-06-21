CREATE TABLE IF NOT EXISTS performance_favorites (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    performance_id uuid NOT NULL REFERENCES performance(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, performance_id)
);

CREATE INDEX IF NOT EXISTS performance_favorites_user_created_idx
    ON performance_favorites (user_id, created_at DESC);
