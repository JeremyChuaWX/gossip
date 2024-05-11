CREATE UNIQUE INDEX IF NOT EXISTS room_users_index ON room_users (
    room_id,
    user_id
)
