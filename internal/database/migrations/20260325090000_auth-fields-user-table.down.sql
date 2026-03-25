ALTER TABLE users
DROP COLUMN IF EXISTS last_login_at,
DROP COLUMN IF EXISTS is_active,
DROP COLUMN IF EXISTS refresh_token_hash,
DROP COLUMN IF EXISTS password_hash;
