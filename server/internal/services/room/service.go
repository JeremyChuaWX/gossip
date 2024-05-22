package room

import (
	"context"
	"gossip/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	PgPool *pgxpool.Pool
}

func (s *Service) Create(
	ctx context.Context,
	dto CreateDTO,
) (models.Room, error) {
	sql := `
	INSERT INTO rooms (
		name
	)
	VALUES (
		$1
	)
	RETURNING
		id,
		name
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql, dto.Name)
	return pgx.CollectExactlyOneRow[models.Room](rows, pgx.RowToStructByName)
}

func (s *Service) FindOne(
	ctx context.Context,
	dto FindOneDTO,
) (models.Room, error) {
	sql := `
	SELECT
		id,
		name
	FROM rooms
	WHERE
		id = $1
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql, dto.Id)
	return pgx.CollectExactlyOneRow[models.Room](rows, pgx.RowToStructByName)
}

func (s *Service) FindMany(ctx context.Context) ([]models.Room, error) {
	sql := `
	SELECT
		id,
		name
	FROM rooms
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql)
	return pgx.CollectRows[models.Room](rows, pgx.RowToStructByName)
}

func (s *Service) Update(
	ctx context.Context,
	dto UpdateDTO,
) (models.Room, error) {
	sql := `
	UPDATE rooms
	SET
		name = COALESCE($1, name)
	WHERE
		id = $2
	RETURNING
		id,
		username,
		password_hash
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql, dto.Name, dto.Id)
	return pgx.CollectExactlyOneRow[models.Room](rows, pgx.RowToStructByName)
}

func (s *Service) Delete(
	ctx context.Context,
	dto DeleteDTO,
) (models.Room, error) {
	sql := `
	DELETE FROM rooms
	WHERE
		id = $1
	RETURNING
		id,
		name
	;
	`
	rows, _ := s.PgPool.Query(ctx, sql, dto.Id)
	return pgx.CollectExactlyOneRow[models.Room](rows, pgx.RowToStructByName)
}
