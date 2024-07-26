package message

import (
	"context"
	"gossip/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	PgPool *pgxpool.Pool
}

func (s *Service) Save(
	ctx context.Context,
	dto SaveDto,
) (models.Message, error) {
	sql := `
	INSERT INTO messages (
		user_id,
		room_id,
		body
	)
	VALUES (
		$1,
		$2,
		$3
	)
	RETURNING (
		id,
		user_id,
		room_id,
		body,
		timestamp
	)
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql, dto.UserId, dto.RoomId, dto.Body)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Message])
}

func (s *Service) FindManyByRoomId(
	ctx context.Context,
	dto FindManyByRoomIdDto,
) ([]models.Message, error) {
	sql := `
	SELECT 
		id,
		user_id,
		room_id,
		body,
		timestamp
	FROM messages
	WHERE
		room_id = $1
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql, dto.RoomId)
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Message])
}
