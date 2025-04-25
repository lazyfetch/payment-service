package storage

import "errors"

var (
	ErrUserIDNotFound = errors.New("user_id not found")

	ErrUserIDNotExists = errors.New("user_id not exists")
)
