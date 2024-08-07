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
	UserId uuid.UUID `db:"id" json:"userId"`
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
	rows, err := r.PgPool.Query(ctx, sql, dto.Username, dto.PasswordHash)
	defer rows.Close()
	if err != nil {
		return UserCreateResult{}, err
	}
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[UserCreateResult],
	)
}

type UserFindOneParams struct {
	UserId uuid.UUID
}

type UserFindOneResult struct {
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
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
	rows, err := r.PgPool.Query(ctx, sql, dto.UserId)
	defer rows.Close()
	if err != nil {
		return UserFindOneResult{}, err
	}
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[UserFindOneResult],
	)
}

type UserFindOneByUsernameParams struct {
	Username string
}

type UserFindOneByUsernameResult struct {
	UserId       uuid.UUID `db:"id" json:"userId"`
	PasswordHash string    `db:"password_hash" json:"-"`
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
	rows, err := r.PgPool.Query(ctx, sql, dto.Username)
	defer rows.Close()
	if err != nil {
		return UserFindOneByUsernameResult{}, err
	}
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[UserFindOneByUsernameResult],
	)
}

type UsersFindManyByRoomIdParams struct {
	RoomId uuid.UUID
}

type UsersFindManyByRoomIdResult struct {
	UserId   uuid.UUID `db:"id" json:"userId"`
	Username string    `db:"username" json:"username"`
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
	rows, err := r.PgPool.Query(ctx, sql, dto.RoomId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
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
	rows, err := r.PgPool.Query(
		ctx,
		sql,
		dto.Username,
		dto.PasswordHash,
		dto.UserId,
	)
	defer rows.Close()
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
	rows, err := r.PgPool.Query(ctx, sql, dto.UserId)
	defer rows.Close()
	return err
}

type UserCheckRoomMembershipParams struct {
	UserId uuid.UUID
	RoomId uuid.UUID
}

type UserCheckRoomMembershipResult struct {
	Membership int `db:"membership"`
}

func (r *Repository) UserCheckRoomMembership(
	ctx context.Context,
	dto UserCheckRoomMembershipParams,
) (bool, error) {
	sql := `
	SELECT
		COUNT(user_id) as membership
	FROM room_users
	WHERE
		1 = 1
		AND user_id = $1
		AND room_id = $2
	;
	`
	rows, err := r.PgPool.Query(ctx, sql, dto.UserId, dto.RoomId)
	defer rows.Close()
	if err != nil {
		return false, err
	}
	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[UserCheckRoomMembershipResult],
	)
	return result.Membership > 0, err
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
	rows, err := r.PgPool.Query(ctx, sql, dto.UserId, dto.RoomId)
	defer rows.Close()
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
	rows, err := r.PgPool.Query(ctx, sql, dto.UserId, dto.RoomId)
	defer rows.Close()
	return err
}

type SessionCreateParams struct {
	UserId uuid.UUID
}

type SessionCreateResult struct {
	SessionId uuid.UUID `db:"id" json:"sessionId"`
	ExpiresOn time.Time `db:"expires_on" json:"expiresOn"`
}

func (r *Repository) SessionCreate(
	ctx context.Context,
	dto SessionCreateParams,
) (SessionCreateResult, error) {
	sql := `
	INSERT INTO user_sessions (
		user_id,
		expires_on
	)
	VALUES (
		$1,
		$2
	)
	RETURNING
		id,
		expires_on
	;
	`
	expiresOn := time.Now().Add(SESSION_DURATION)
	rows, err := r.PgPool.Query(ctx, sql, dto.UserId, expiresOn)
	defer rows.Close()
	if err != nil {
		return SessionCreateResult{}, err
	}
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[SessionCreateResult],
	)
}

type SessionFindOneParams struct {
	SessionId uuid.UUID
}

type SessionFindOneResult struct {
	SessionId uuid.UUID `db:"id" json:"sessionId"`
	UserId    uuid.UUID `db:"user_id" json:"userId"`
	Username  string    `db:"username" json:"username"`
	ExpiresOn time.Time `db:"expires_on" json:"expiresOn"`
}

func (r *Repository) SessionFindOne(
	ctx context.Context,
	dto SessionFindOneParams,
) (SessionFindOneResult, error) {
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
	rows, err := r.PgPool.Query(ctx, sql, dto.SessionId)
	defer rows.Close()
	if err != nil {
		return SessionFindOneResult{}, err
	}
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[SessionFindOneResult],
	)
}

type SessionDeleteParams struct {
	SessionId uuid.UUID
}

func (r *Repository) SessionDelete(
	ctx context.Context,
	dto SessionDeleteParams,
) error {
	sql := `
	DELETE FROM user_sessions
	WHERE
		id = $1
	;
	`
	rows, err := r.PgPool.Query(ctx, sql, dto.SessionId)
	defer rows.Close()
	return err
}

type RoomCreateParams struct {
	Name string
}

type RoomCreateResult struct {
	RoomId uuid.UUID `db:"id" json:"roomId"`
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
	rows, err := r.PgPool.Query(ctx, sql, dto.Name)
	defer rows.Close()
	if err != nil {
		return RoomCreateResult{}, err
	}
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[RoomCreateResult],
	)
}

type RoomFindOneParams struct {
	RoomId uuid.UUID
}

type RoomFindOneResult struct {
	Name string `db:"name" json:"name"`
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
	rows, err := r.PgPool.Query(ctx, sql, dto.RoomId)
	defer rows.Close()
	if err != nil {
		return RoomFindOneResult{}, err
	}
	return pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[RoomFindOneResult],
	)
}

type RoomFindManyResult struct {
	RoomId uuid.UUID `db:"id" json:"roomId"`
	Name   string    `db:"name" json:"name"`
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
	rows, err := r.PgPool.Query(ctx, sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[RoomFindManyResult])
}

type RoomFindManyByUserIdParams struct {
	UserId uuid.UUID
}

type RoomFindManyByUserIdResult struct {
	RoomId uuid.UUID `db:"id" json:"roomId"`
	Name   string    `db:"name" json:"name"`
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
	rows, err := r.PgPool.Query(ctx, sql, dto.UserId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
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
	rows, err := r.PgPool.Query(ctx, sql, dto.Name, dto.RoomId)
	defer rows.Close()
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
	rows, err := r.PgPool.Query(ctx, sql, dto.RoomId)
	defer rows.Close()
	return err
}

type MessageSaveParams struct {
	UserId uuid.UUID
	RoomId uuid.UUID
	Body   string
}

func (r *Repository) MessageSave(
	ctx context.Context,
	dto MessageSaveParams,
) error {
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
	;
	`
	rows, err := r.PgPool.Query(ctx, sql, dto.UserId, dto.RoomId, dto.Body)
	defer rows.Close()
	return err
}

type MessagesFindManyByRoomIdParams struct {
	RoomId uuid.UUID
}

type MessagesFindManyByRoomIdResult struct {
	MessageId uuid.UUID `db:"id" json:"messageId"`
	UserId    uuid.UUID `db:"user_id" json:"userId"`
	RoomId    uuid.UUID `db:"room_id" json:"roomId"`
	Username  string    `db:"username" json:"username"`
	Body      string    `db:"body" json:"body"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}

func (r *Repository) MessagesFindManyByRoomId(
	ctx context.Context,
	dto MessagesFindManyByRoomIdParams,
) ([]MessagesFindManyByRoomIdResult, error) {
	sql := `
	SELECT
		messages.id,
		messages.user_id,
		messages.room_id,
		users.username,
		messages.body,
		messages.timestamp
	FROM messages
		INNER JOIN users ON users.id = messages.user_id
	WHERE
		room_id = $1
	ORDER BY
		messages.timestamp ASC
	;
	`
	rows, err := r.PgPool.Query(ctx, sql, dto.RoomId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(
		rows,
		pgx.RowToStructByName[MessagesFindManyByRoomIdResult],
	)
}
