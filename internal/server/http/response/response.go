package response

import "github.com/stuton/xm-golang-exercise/internal/errors"

type Response struct {
	Data   interface{}     `json:"data,omitempty"`
	Errors []*errors.Error `json:"errors,omitempty"`
}

func NewSuccessResponse(data interface{}) *Response {
	res := &Response{
		Data: data,
	}
	return res
}

func NewErrResponse(errs ...error) *Response {

	result := make([]*errors.Error, 0)
	for _, err := range errs {
		if e, ok := err.(*errors.Error); ok {
			result = append(result, e)
		}
	}

	res := &Response{
		Errors: result,
	}

	return res
}
