package user

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
) (User, error) {
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
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) FindOne(
	ctx context.Context,
	dto FindOneDTO,
) (User, error) {
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
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) FindOneByUsername(
	ctx context.Context,
	dto FindOneByUsernameDTO,
) (User, error) {
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
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) Update(
	ctx context.Context,
	dto UpdateDTO,
) (User, error) {
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
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) Delete(
	ctx context.Context,
	dto DeleteDTO,
) (User, error) {
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
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}
