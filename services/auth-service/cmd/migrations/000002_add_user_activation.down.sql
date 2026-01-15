DROP TABLE IF EXISTS user_activation_tokens;
ALTER TABLE users DROP COLUMN IF EXISTS is_active;