package roomuser

import (
	"context"
	"gossip/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	PgPool *pgxpool.Pool
}

func (s *Service) UserJoinRoom(
	ctx context.Context,
	dto UserJoinRoomDTO,
) (models.RoomUser, error) {
	sql := `
	INSERT INTO roomusers (
		user_id,
		room_id
	)
	VALUES (
		$1,
		$2
	)
	RETURNING
		user_id,
		room_id
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql, dto.UserId, dto.RoomId)
	return pgx.CollectExactlyOneRow[models.RoomUser](
		rows,
		pgx.RowToStructByName,
	)
}

func (s *Service) UserLeaveRoom(
	ctx context.Context,
	dto UserLeaveRoomDTO,
) (models.RoomUser, error) {
	sql := `
	DELETE FROM roomusers
	WHERE
		1 = 1
		AND user_id = $1
		AND room_id = $2
	RETURNING
		user_id,
		room_id
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql, dto.UserId, dto.RoomId)
	return pgx.CollectExactlyOneRow[models.RoomUser](
		rows,
		pgx.RowToStructByName,
	)
}

func (s *Service) FindRoomIdsByUserId(
	ctx context.Context,
	dto FindRoomIdsByUserIdDTO,
) ([]models.RoomUser, error) {
	sql := `
	SELECT
		user_id,
		room_id
	FROM roomusers
	WHERE
		user_id = $1
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql, dto.UserId)
	return pgx.CollectRows[models.RoomUser](rows, pgx.RowToStructByName)
}

func (s *Service) FindUserIdsByRoomId(
	ctx context.Context,
	dto FindUserIdsByRoomIdDTO,
) ([]models.RoomUser, error) {
	sql := `
	SELECT
		user_id,
		room_id
	FROM roomusers
	WHERE
		room_id = $1
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql, dto.RoomId)
	return pgx.CollectRows[models.RoomUser](rows, pgx.RowToStructByName)
}
