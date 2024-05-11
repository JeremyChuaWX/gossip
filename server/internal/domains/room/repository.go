package room

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
) (Room, error) {
	sql := `
	INSERT INTO rooms (
		name
	) VALUES (
		$1
	) RETURNING
		id,
		name
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Name)
	return pgx.CollectExactlyOneRow[Room](rows, pgx.RowToStructByName)
}

func (r *Repository) FindOne(
	ctx context.Context,
	dto FindOneDTO,
) (Room, error) {
	sql := `
	SELECT
		id,
		name
	FROM rooms WHERE
		id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Id)
	return pgx.CollectExactlyOneRow[Room](rows, pgx.RowToStructByName)
}

func (r *Repository) Update(
	ctx context.Context,
	dto UpdateDTO,
) (Room, error) {
	sql := `
	UPDATE rooms SET
		name = COALESCE($1, name)
	WHERE
		id = $2
	RETURNING
		id,
		username,
		password_hash
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Name, dto.Id)
	return pgx.CollectExactlyOneRow[Room](rows, pgx.RowToStructByName)
}

func (r *Repository) Delete(
	ctx context.Context,
	dto DeleteDTO,
) (Room, error) {
	sql := `
	DELETE FROM rooms WHERE
		id = $1
	RETURNING
		id,
		name
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Id)
	return pgx.CollectExactlyOneRow[Room](rows, pgx.RowToStructByName)
}
