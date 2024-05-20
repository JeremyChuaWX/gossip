CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID REFERENCES room(id),
    user_id UUID REFERENCES user(id),
    body TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
