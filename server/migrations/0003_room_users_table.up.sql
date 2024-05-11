CREATE TABLE IF NOT EXISTS room_users (
    room_id UUID REFERENCES room(id),
    user_id UUID REFERENCES user(id)
)
