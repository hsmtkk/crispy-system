package sessionstore

import "fmt"

type RedisImpl struct{}

func NewRedisImpl(redisHost string, redisPort int) (SessionStore, error) {
	return &RedisImpl{}, nil
}

func (r *RedisImpl) NewSession(userID string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (r *RedisImpl) GetUserID(sessionID string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (r *RedisImpl) DeleteSession(sessionID string) error {
	return fmt.Errorf("not implemented")
}
