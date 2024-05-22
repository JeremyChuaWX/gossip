CREATE TABLE IF NOT EXISTS room_users (
    room_id UUID NOT NULL REFERENCES rooms(id),
    user_id UUID NOT NULL REFERENCES users(id),
    PRIMARY KEY (room_id, user_id)
)
