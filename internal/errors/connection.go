package errors

import "errors"

var (
	WrongPort error = errors.New("wrong port")
)
