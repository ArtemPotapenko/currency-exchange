package error

import (
	"fmt"
)

type NotFoundError struct {
	Message string
	Detail  string
}

func (e *NotFoundError) Error() string {
	if e.Detail == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Detail)
}

func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

type ValidationError struct {
	Message string
	Detail  string
}

func (e *ValidationError) Error() string {
	if e.Detail == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Detail)
}

func (e *ValidationError) Is(target error) bool {
	_, ok := target.(*ValidationError)
	return ok
}

type InternalError struct {
	Message string
	Detail  string
}

func (e *InternalError) Error() string {
	if e.Detail == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Detail)
}

func (e *InternalError) Is(target error) bool {
	_, ok := target.(*InternalError)
	return ok
}

var (
	ErrNotFound   = &NotFoundError{}
	ErrValidation = &ValidationError{}
	ErrInternal   = &InternalError{}
)

func NotFound(message string, detail string) error {
	return &NotFoundError{Message: message, Detail: detail}
}

func Validation(message string, detail string) error {
	return &ValidationError{Message: message, Detail: detail}
}

func Internal(message string, detail string) error {
	return &InternalError{Message: message, Detail: detail}
}
