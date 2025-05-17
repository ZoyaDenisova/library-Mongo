package errors

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserBlocked         = errors.New("user is blocked")
	ErrInvalidID           = errors.New("invalid id")
	ErrBookNotFound        = errors.New("resource not found")
	ErrBookAlreadyBorrowed = errors.New("book is already borrowed")
	ErrBorrowNotFound      = errors.New("borrow not found")
	ErrAlreadyReturned     = errors.New("book already returned")
)
