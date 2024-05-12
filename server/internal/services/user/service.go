package user

import (
	"context"
	"gossip/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	PgPool *pgxpool.Pool
}

func (r *Service) Create(
	ctx context.Context,
	dto CreateDTO,
) (models.User, error) {
	sql := `
	INSERT INTO users (
		username,
		password_hash
	) VALUES (
		$1,
		$2
	) RETURNING
		id,
		username,
		password_hash
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Username, dto.PasswordHash)
	return pgx.CollectExactlyOneRow[models.User](rows, pgx.RowToStructByName)
}

func (r *Service) FindOne(
	ctx context.Context,
	dto FindOneDTO,
) (models.User, error) {
	sql := `
	SELECT
		id,
		username,
		password_hash
	FROM users WHERE
		id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Id)
	return pgx.CollectExactlyOneRow[models.User](rows, pgx.RowToStructByName)
}

func (r *Service) FindOneByUsername(
	ctx context.Context,
	dto FindOneByUsernameDTO,
) (models.User, error) {
	sql := `
	SELECT
		id,
		username,
		password_hash
	FROM users WHERE
		username = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Username)
	return pgx.CollectExactlyOneRow[models.User](rows, pgx.RowToStructByName)
}

func (r *Service) Update(
	ctx context.Context,
	dto UpdateDTO,
) (models.User, error) {
	sql := `
	UPDATE users SET
		username = COALESCE($1, username),
		password_hash = COALESCE($2, password_hash)
	WHERE
		id = $3
	RETURNING
		id,
		username,
		password_hash
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Username, dto.PasswordHash, dto.Id)
	return pgx.CollectExactlyOneRow[models.User](rows, pgx.RowToStructByName)
}

func (r *Service) Delete(
	ctx context.Context,
	dto DeleteDTO,
) (models.User, error) {
	sql := `
	DELETE FROM users WHERE
		id = $1
	RETURNING
		id,
		username,
		password_hash
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Id)
	return pgx.CollectExactlyOneRow[models.User](rows, pgx.RowToStructByName)
}
