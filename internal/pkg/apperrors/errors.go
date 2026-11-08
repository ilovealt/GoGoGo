package apperrors

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrInvalid      = errors.New("invalid argument")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrTokenRevoked = errors.New("refresh token revoked")
	ErrTokenExpired = errors.New("refresh token expired")
)

// Error 携带业务分类和可直接返回给调用方的消息。
// errors.Is(err, ErrNotFound) 等仍然可用。
type Error struct {
	Kind    error
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Kind
}

func New(kind error, message string) error {
	return &Error{Kind: kind, Message: message}
}
