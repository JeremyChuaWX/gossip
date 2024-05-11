package roomuser

import (
	"context"

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
	) VALUES (
		$1,
		$2
	) RETURNING
		room_id,
		user_id
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.RoomId, dto.UserId)
	return pgx.CollectExactlyOneRow[RoomUser](rows, pgx.RowToStructByName)
}

func (r *Repository) FindManyByRoomId(
	ctx context.Context,
	dto FindManyByRoomIdDTO,
) (RoomUser, error) {
	sql := `
	SELECT
		room_id,
		user_id
	FROM room_users WHERE
		room_id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.RoomId)
	return pgx.CollectExactlyOneRow[RoomUser](rows, pgx.RowToStructByName)
}

func (r *Repository) FindManyByUserId(
	ctx context.Context,
	dto FindManyByUserIdDTO,
) (RoomUser, error) {
	sql := `
	SELECT
		room_id,
		user_id
	FROM room_users WHERE
		user_id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.UserId)
	return pgx.CollectExactlyOneRow[RoomUser](rows, pgx.RowToStructByName)
}

func (r *Repository) Delete(
	ctx context.Context,
	dto DeleteDTO,
) (RoomUser, error) {
	sql := `
	DELETE FROM room_users WHERE
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
