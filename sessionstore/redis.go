package sessionstore

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisImpl struct {
	client *redis.Client
}

func NewRedisImpl(redisHost string, redisPort int) SessionStore {
	addr := fmt.Sprintf("%s:%d", redisHost, redisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &RedisImpl{client}
}

func (r *RedisImpl) NewSession(ctx context.Context, sessionID, userID string) (string, error) {
	if err := r.client.Set(ctx, sessionID, userID, 0).Err(); err != nil {
		return "", fmt.Errorf("redis.Client.Set failed; %w", err)
	}
	return sessionID, nil
}

func (r *RedisImpl) GetUserID(ctx context.Context, sessionID string) (string, error) {
	userID, err := r.client.Get(ctx, sessionID).Result()
	if err != nil {
		return "", fmt.Errorf("redis.Client.Get failed; %w", err)
	}
	return userID, nil
}

func (r *RedisImpl) DeleteSession(ctx context.Context, sessionID string) error {
	if err := r.client.Del(ctx, sessionID).Err(); err != nil {
		return fmt.Errorf("redis.Client.Del failed; %w", err)
	}
	return nil
}
