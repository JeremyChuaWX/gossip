package session

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/redis/go-redis/v9"
)

const SESSION_EXPIRATION = time.Hour * 24 * 7

type Repository struct {
	Redis *redis.Client
}

func (r *Repository) Create(
	ctx context.Context,
	userId string,
) (string, error) {
	sessionId, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	if err = r.Redis.Set(ctx, sessionId.String(), userId, SESSION_EXPIRATION).Err(); err != nil {
		return "", err
	}
	return sessionId.String(), nil
}

func (r *Repository) Get(
	ctx context.Context,
	sessionId string,
) (string, error) {
	return r.Redis.Get(ctx, sessionId).Result()
}

func (r *Repository) Delete(
	ctx context.Context,
	sessionId string,
) error {
	return r.Redis.Del(ctx, sessionId).Err()
}
