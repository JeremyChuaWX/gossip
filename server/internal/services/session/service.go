package session

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Redis *redis.Client
}

func (s *Service) Create(
	ctx context.Context,
	userId uuid.UUID,
) (string, error) {
	sessionId, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	if err = s.Redis.Set(ctx, sessionId.String(), userId.String(), SESSION_EXPIRATION).Err(); err != nil {
		return "", err
	}
	return sessionId.String(), nil
}

func (s *Service) Get(
	ctx context.Context,
	sessionId string,
) (uuid.UUID, error) {
	userIdString, err := s.Redis.Get(ctx, sessionId).Result()
	if err != nil {
		return uuid.Nil, err
	}
	userId, err := uuid.FromString(userIdString)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, err
}

func (s *Service) Delete(
	ctx context.Context,
	sessionId string,
) error {
	return s.Redis.Del(ctx, sessionId).Err()
}
