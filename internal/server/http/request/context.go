package request

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stuton/xm-golang-exercise/internal/errors"
	"github.com/stuton/xm-golang-exercise/internal/server/http/response"
)

type context struct {
	*gin.Context
}

func FromContext(c *gin.Context) *context {
	return &context{
		Context: c,
	}
}

func (s *context) Ok(data interface{}) {
	s.JSON(http.StatusOK, response.NewSuccessResponse(data))
}

func (s *context) Err(statusCode int, errs ...error) {
	s.JSON(statusCode, response.NewErrResponse(errs...))
}

func (s *context) BadRequest(errs ...error) {
	s.Err(http.StatusBadRequest, errs...)
}

func (s *context) ResponseError(errs ...error) {
	statusCode := http.StatusInternalServerError

	errorList := make([]error, 0)

	for _, err := range errs {
		if e, ok := err.(*errors.Error); ok {
			statusCode = e.StatusCode
			errorList = append(errorList, e)
			continue
		}
		errorList = append(errorList, err)
	}

	s.JSON(statusCode, response.NewErrResponse(errorList...))
}
