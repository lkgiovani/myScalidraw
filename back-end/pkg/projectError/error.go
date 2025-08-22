package projectError

import (
	"errors"
	"fmt"
)

const (
	ECONFLICT       = "conflict"
	EINTERNAL       = "internal"
	EINVALID        = "invalid"
	ENOTFOUND       = "not_found"
	ENOTIMPLEMENTED = "not_implemented"
	EUNAUTHORIZED   = "unauthorized"
)

type Error struct {
	Code      string
	Message   string
	Path      string
	PrevError error
}

func (e *Error) Error() string {
	if e.PrevError != nil {
		return fmt.Sprintf("error: code=%s message=%s prevError=%v", e.Code, e.Message, e.PrevError)
	}
	return fmt.Sprintf("error: code=%s message=%s", e.Code, e.Message)
}

func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	}

	if errors.As(err, &e) {
		return e.Code
	}

	return EINTERNAL
}

func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error."
}

func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
