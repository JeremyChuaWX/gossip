package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var redisInitialisationError = errors.New("error initialising redis")

// remember to defer close the connection
func Init(
	ctx context.Context,
	addr string,
	password string,
) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	if rdb == nil {
		return nil, redisInitialisationError
	}
	return rdb, nil
}
