ALTER TABLE users
ADD COLUMN IF NOT EXISTS is_active boolean NOT NULL DEFAULT false;

CREATE TABLE IF NOT EXISTS user_activation_tokens (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token text NOT NULL UNIQUE,
  expiry timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_user_activation_tokens_user_id
ON user_activation_tokens(user_id);