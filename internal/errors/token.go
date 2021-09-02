package errors

import "errors"

var (
	TockenExpired  error = errors.New("tocken expired")
	InvalidTocken  error = errors.New("invalid tocken")
	JWTKeyNotFound error = errors.New("key not found")
)
