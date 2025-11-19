package user

import "errors"

var ErrEmptyUsername = errors.New("username cannot be empty")
var ErrEmptyPasswordHash = errors.New("password hash cannot be empty")
