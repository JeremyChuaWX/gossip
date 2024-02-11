package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	PgPool *pgxpool.Pool
}

func (r *Repository) create(
	ctx context.Context,
	dto createDTO,
) (User, error) {
	sql := `
	INSERT INTO users (
		username,
		password
	) VALUES (
		$1,
		$2
	) RETURNING (
		id,
		username,
		password
	);
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Username, dto.Password)
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) findOne(
	ctx context.Context,
	dto findOneDTO,
) (User, error) {
	sql := `
	SELECT (
		id,
		username,
		password
	) FROM users WHERE (
		id = $1,
	);
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Id)
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) update(
	ctx context.Context,
	dto updateDTO,
) (User, error) {
	sql := `
	UPDATE users SET (
		username = COALESCE($1, username),
		password = COALESCE($2, password)
	) WHERE (
		id = $3,
	) RETURNING (
		id,
		username,
		password
	);
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Username, dto.Password, dto.Id)
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}

func (r *Repository) delete(
	ctx context.Context,
	dto deleteDTO,
) (User, error) {
	sql := `
	DELETE FROM users WHERE (
		id = $1,
	) RETURNING (
		id,
		username,
		password
	);
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Id)
	return pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName)
}
