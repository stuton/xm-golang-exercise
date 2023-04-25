package errors

import "net/http"

type Error struct {
	Key        string      `json:"key,omitempty"`
	Value      interface{} `json:"value,omitempty"`
	Message    string      `json:"message,omitempty"`
	StatusCode int         `json:"-"`
}

func NewError(key string, value interface{}, message string, statusCode int) *Error {
	return &Error{Key: key, Value: value, Message: message, StatusCode: statusCode}
}

func NewBadRequestError(key string, value interface{}, message string) *Error {
	return &Error{Key: key, Value: value, Message: message, StatusCode: http.StatusBadRequest}
}

func NewNotFoundError(key string, value interface{}, message string) *Error {
	return &Error{Key: key, Value: value, Message: message, StatusCode: http.StatusNotFound}
}

func NewConflictError(key string, value interface{}, message string) *Error {
	return &Error{Key: key, Value: value, Message: message, StatusCode: http.StatusConflict}
}

func NewUnauthorizedError(key string, value interface{}, message string) *Error {
	return &Error{Key: key, Value: value, Message: message, StatusCode: http.StatusUnauthorized}
}

func NewForbiddenError(key string, value interface{}, message string) *Error {
	return &Error{Key: key, Value: value, Message: message, StatusCode: http.StatusForbidden}
}

func NewInternalServerError(key string, value interface{}, message string) *Error {
	return &Error{Key: key, Value: value, Message: message, StatusCode: http.StatusInternalServerError}
}

func NewNoContentError(key string, value interface{}, message string) *Error {
	return &Error{Key: key, Value: value, Message: message, StatusCode: http.StatusNoContent}
}

func (s *Error) Error() string {
	return s.Message
}

func (s *Error) WithMessage(message string) *Error {
	s.Message = message
	return s
}

func (s *Error) WithKey(key string) *Error {
	s.Key = key
	return s
}

func (s *Error) WithValue(value interface{}) *Error {
	s.Value = value
	return s
}

func (s *Error) WithStatusCode(statusCode int) *Error {
	s.StatusCode = statusCode
	return s
}
