package auth

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/redis/go-redis/v9"
)

const SESSION_EXPIRATION = time.Hour * 24 * 7

type Repository struct {
	redis *redis.Client
}

func (r *Repository) createSession(
	ctx context.Context,
	userId string,
) (string, error) {
	sessionId, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	if err = r.redis.Set(ctx, sessionId.String(), userId, SESSION_EXPIRATION).Err(); err != nil {
		return "", err
	}
	return sessionId.String(), nil
}

func (r *Repository) getSession(
	ctx context.Context,
	sessionId string,
) (string, error) {
	return r.redis.Get(ctx, sessionId).Result()
}

func (r *Repository) deleteSession(
	ctx context.Context,
	sessionId string,
) error {
	return r.redis.Del(ctx, sessionId).Err()
}
