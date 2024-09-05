package data

import "errors"

var (
	ErrNotFound = errors.New("data: not found")
	errInternal = errors.New("data: internal error")
)

func ErrInternal(err error) error {
	return errors.Join(errInternal, err)
}
