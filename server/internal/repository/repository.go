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
	UserId       uuid.UUID `db:"id"`
	PasswordHash string    `db:"password_hash"`
}

func (r *Repository) UserFindOneByUsername(
	ctx context.Context,
	dto UserFindOneByUsernameParams,
) (UserFindOneByUsernameResult, error) {
	sql := `
	SELECT
		id,
		password_hash
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

type UsersFindManyByRoomIdParams struct {
	RoomId uuid.UUID
}

type UsersFindManyByRoomIdResult struct {
	UserId   uuid.UUID `db:"id"`
	Username string    `db:"username"`
}

func (r *Repository) UsersFindManyByRoomId(
	ctx context.Context,
	dto UsersFindManyByRoomIdParams,
) ([]UsersFindManyByRoomIdResult, error) {
	sql := `
	SELECT
		users.id,
		users.username
	FROM room_users
		INNER JOIN users ON users.id = room_users.user_id
	WHERE
		room_users.room_id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.RoomId)
	return pgx.CollectRows(
		rows,
		pgx.RowToStructByName[UsersFindManyByRoomIdResult],
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

type UserJoinRoomParams struct {
	UserId uuid.UUID
	RoomId uuid.UUID
}

func (r *Repository) UserJoinRoom(
	ctx context.Context,
	dto UserJoinRoomParams,
) error {
	sql := `
	INSERT INTO room_users (
		user_id,
		room_id
	)
	VALUES (
		$1,
		$2
	)
	;
	`
	_, err := r.PgPool.Query(ctx, sql, dto.UserId, dto.RoomId)
	return err
}

type UserLeaveRoomParams struct {
	UserId uuid.UUID
	RoomId uuid.UUID
}

func (r *Repository) UserLeaveRoom(
	ctx context.Context,
	dto UserLeaveRoomParams,
) error {
	sql := `
	DELETE FROM room_users
	WHERE
		1 = 1
		AND user_id = $1
		AND room_id = $2
	;
	`
	_, err := r.PgPool.Query(ctx, sql, dto.UserId, dto.RoomId)
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

type UserSessionDeleteParams struct {
	SessionId uuid.UUID
}

func (r *Repository) UserSessionDelete(
	ctx context.Context,
	dto UserSessionDeleteParams,
) error {
	sql := `
	DELETE FROM user_sessions
	WHERE
		id = $1
	;
	`
	_, err := r.PgPool.Query(ctx, sql, dto.SessionId)
	return err
}

type RoomCreateParams struct {
	Name string
}

type RoomCreateResult struct {
	RoomId uuid.UUID `db:"id"`
}

func (r *Repository) RoomCreate(
	ctx context.Context,
	dto RoomCreateParams,
) (RoomCreateResult, error) {
	sql := `
	INSERT INTO rooms (
		name
	)
	VALUES (
		$1
	)
	RETURNING
		id
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.Name)
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[RoomCreateResult],
	)
}

type RoomFindOneParams struct {
	RoomId uuid.UUID
}

type RoomFindOneResult struct {
	Name string `db:"name"`
}

func (r *Repository) RoomFindOne(
	ctx context.Context,
	dto RoomFindOneParams,
) (RoomFindOneResult, error) {
	sql := `
	SELECT
		name
	FROM rooms
	WHERE
		id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.RoomId)
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[RoomFindOneResult],
	)
}

type RoomFindManyResult struct {
	RoomId uuid.UUID `db:"id"`
	Name   string    `db:"name"`
}

func (r *Repository) RoomFindMany(
	ctx context.Context,
) ([]RoomFindManyResult, error) {
	sql := `
	SELECT
		id,
		name
	FROM rooms
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql)
	return pgx.CollectRows(rows, pgx.RowToStructByName[RoomFindManyResult])
}

type RoomFindManyByUserIdParams struct {
	UserId uuid.UUID
}

type RoomFindManyByUserIdResult struct {
	RoomId uuid.UUID `db:"id"`
	Name   string    `db:"name"`
}

func (r *Repository) RoomFindManyByUserId(
	ctx context.Context,
	dto RoomFindManyByUserIdParams,
) ([]RoomFindManyByUserIdResult, error) {
	sql := `
	SELECT
		rooms.id,
		rooms.name
	FROM room_users
		INNER JOIN rooms ON rooms.id = room_users.room_id
	WHERE
		room_users.user_id = $1
	;
	`
	rows, _ := r.PgPool.Query(ctx, sql, dto.UserId)
	return pgx.CollectRows(
		rows,
		pgx.RowToStructByName[RoomFindManyByUserIdResult],
	)
}

type RoomUpdateParams struct {
	RoomId uuid.UUID
	Name   *string
}

func (r *Repository) RoomUpdate(
	ctx context.Context,
	dto RoomUpdateParams,
) error {
	sql := `
	UPDATE rooms
	SET
		name = COALESCE($1, name)
	WHERE
		id = $2
	;
	`
	_, err := r.PgPool.Query(ctx, sql, dto.Name, dto.RoomId)
	return err
}

type RoomDeleteParams struct {
	RoomId uuid.UUID
}

func (r *Repository) RoomDelete(
	ctx context.Context,
	dto RoomDeleteParams,
) error {
	sql := `
	DELETE FROM rooms
	WHERE
		id = $1
	;
	`
	_, err := r.PgPool.Query(ctx, sql, dto.RoomId)
	return err
}
