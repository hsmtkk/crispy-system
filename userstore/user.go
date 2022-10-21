package userstore

import "context"

type UserStore interface {
	Increment(ctx context.Context, userID string) (int, error)
}
