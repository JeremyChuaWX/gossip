CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID REFERENCES rooms(id),
    user_id UUID REFERENCES users(id),
    body TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
