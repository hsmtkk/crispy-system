package userstore

type UserStore interface {
	Increment(userID string) (int, error)
}
