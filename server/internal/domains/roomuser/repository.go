package roomuser

import (
	"context"
	"gossip/internal/domains/room"
	"gossip/internal/domains/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	PgPool *pgxpool.Pool
}

func (r *Repository) Create(
	ctx context.Context,
	dto CreateDTO,
) (RoomUser, error) {
	sql := `
	INSERT INTO room_users (
		room_id,
		user_id
	)
	VALUES (
		$1,
		$2
	)
	RETURNING
		room_id,
		user_id
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.RoomId, dto.UserId)
	return pgx.CollectExactlyOneRow[RoomUser](rows, pgx.RowToStructByName)
}

func (r *Repository) Delete(
	ctx context.Context,
	dto DeleteDTO,
) (RoomUser, error) {
	sql := `
	DELETE FROM room_users
	WHERE
		room_id = $1,
		user_id = $2
	RETURNING
		room_id,
		user_id
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.RoomId, dto.UserId)
	return pgx.CollectExactlyOneRow[RoomUser](rows, pgx.RowToStructByName)
}

func (r *Repository) FindManyUsersByRoomId(
	ctx context.Context,
	dto FindManyByRoomIdDTO,
) ([]user.User, error) {
	sql := `
	SELECT
		users.id,
		users.username,
		users.password_hash
	FROM room_users
		INNER JOIN users ON users.id = room_users.user_id
	WHERE
		room_users.room_id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.RoomId)
	return pgx.CollectRows[user.User](rows, pgx.RowToStructByName)
}

func (r *Repository) FindManyRoomsByUserId(
	ctx context.Context,
	dto FindManyByUserIdDTO,
) ([]room.Room, error) {
	sql := `
	SELECT
		rooms.id,
		rooms.name
	FROM room_users
		INNER JOIN rooms ON rooms.id = room_users.room_id
	WHERE
		room_users.user_id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.UserId)
	return pgx.CollectRows[room.Room](rows, pgx.RowToStructByName)
}
