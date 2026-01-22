package courier

import "errors"

var (
	ErrCourierNotFound = errors.New("courier not found")
	ErrCourierExists   = errors.New("courier with this phone already exists")
)
