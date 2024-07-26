CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    expires_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
