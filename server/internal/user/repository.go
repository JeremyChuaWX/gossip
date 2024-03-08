package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	PgPool *pgxpool.Pool
}

func (r *Repository) userCreate(
	ctx context.Context,
	dto userCreateDTO,
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
	rows, _ := r.PgPool.Query(ctx, sql, dto.username, dto.passwordHash)
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) userFindOne(
	ctx context.Context,
	dto userFindOneDTO,
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
	rows, _ := r.PgPool.Query(ctx, sql, dto.id)
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) userFindOneByUsername(
	ctx context.Context,
	dto userFindOneByUsernameDTO,
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
	rows, _ := r.PgPool.Query(ctx, sql, dto.username)
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) userUpdate(
	ctx context.Context,
	dto updateDTO,
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
	rows, _ := r.PgPool.Query(ctx, sql, dto.username, dto.passwordHash, dto.id)
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) userDelete(
	ctx context.Context,
	dto deleteDTO,
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
	rows, _ := r.PgPool.Query(ctx, sql, dto.id)
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}
