package repository

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	PgPool *pgxpool.Pool
}

type UserCreateParams struct {
	Username     string
	PasswordHash string
}

type UserCreateResult struct {
	UserId uuid.UUID `db:"id"`
}

func (r *Repository) UserCreate(
	ctx context.Context,
	dto UserCreateParams,
) (UserCreateResult, error) {
	sql := `
	INSERT INTO users (
		username,
		password_hash
	)
	VALUES (
		$1,
		$2
	)
	RETURNING
		id
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Username, dto.PasswordHash)
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[UserCreateResult],
	)
}

type UserFindOneParams struct {
	UserId uuid.UUID
}

type UserFindOneResult struct {
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}

func (r *Repository) UserFindOne(
	ctx context.Context,
	dto UserFindOneParams,
) (UserFindOneResult, error) {
	sql := `
	SELECT
		username,
		password_hash
	FROM users
	WHERE
		id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.UserId)
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[UserFindOneResult],
	)
}

type UserFindOneByUsernameParams struct {
	Username string
}

type UserFindOneByUsernameResult struct {
	UserId uuid.UUID `db:"id"`
}

func (r *Repository) UserFindOneByUsername(
	ctx context.Context,
	dto UserFindOneByUsernameParams,
) (UserFindOneByUsernameResult, error) {
	sql := `
	SELECT
		id,
	FROM users
	WHERE
		username = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Username)
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[UserFindOneByUsernameResult],
	)
}

type UserUpdateParams struct {
	UserId       uuid.UUID
	Username     *string
	PasswordHash *string
}

func (r *Repository) UserUpdate(
	ctx context.Context,
	dto UserUpdateParams,
) error {
	sql := `
	UPDATE users
	SET
		username = COALESCE($1, username),
		password_hash = COALESCE($2, password_hash)
	WHERE
		id = $3
	;
	`
	_, err := r.PgPool.Query(
		ctx,
		sql,
		dto.Username,
		dto.PasswordHash,
		dto.UserId,
	)
	return err
}

type UserDeleteParams struct {
	UserId uuid.UUID
}

func (r *Repository) UserDelete(
	ctx context.Context,
	dto UserDeleteParams,
) error {
	sql := `
	DELETE FROM users
	WHERE
		id = $1
	;
	`
	_, err := r.PgPool.Query(ctx, sql, dto.UserId)
	return err
}

type UserSessionCreateParams struct {
	UserId uuid.UUID
}

type UserSessionCreateResult struct {
	SessionId uuid.UUID `db:"id"`
	ExpiresOn time.Time `db:"expires_on"`
}

func (r *Repository) UserSessionCreate(
	ctx context.Context,
	dto UserSessionCreateParams,
) (UserSessionCreateResult, error) {
	sql := `
	INSERT INTO user_sessions (
		user_id
	)
	VALUES (
		$1
	)
	RETURNING
		id,
		expires_on
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.UserId)
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[UserSessionCreateResult],
	)
}

type UserSessionFindOneParams struct {
	SessionId uuid.UUID
}

type UserSessionFindOneResult struct {
	UserId    uuid.UUID `db:"user_id"`
	Username  string    `db:"username"`
	ExpiresOn time.Time `db:"expires_on"`
}

func (r *Repository) UserSessionFindOne(
	ctx context.Context,
	dto UserSessionFindOneParams,
) (UserSessionFindOneResult, error) {
	sql := `
	SELECT
		user_sessions.id,
		user_sessions.user_id,
		users.username,
		user_sessions.expires_on
	FROM user_sessions
		INNER JOIN users ON users.id = user_sessions.user_id
	WHERE
		user_sessions.id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.SessionId)
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[UserSessionFindOneResult],
	)
}