package statuserr

import (
	"fmt"
	"github.com/ginx-contribs/ginx/constant/status"
)

// Error represents http response error, which is along with http status code,
// it used to decide how to show error message in response.
type Error struct {
	Err error
	// http status code
	Status status.Status
	// custom error code
	Code int
}

func (e Error) SetCode(code int) Error {
	e.Code = code
	return e
}

func (e Error) SetErrorf(msg string, args ...any) Error {
	e.Err = fmt.Errorf(msg, args...)
	return e
}

func (e Error) SetError(err error) Error {
	e.Err = err
	return e
}

func (e Error) SetStatus(status status.Status) Error {
	e.Status = status
	return e
}

func (e Error) Error() string {
	return e.Err.Error()
}
func New() Error {
	return Error{}
}

// helper functions

func Errorf(msg string, args ...any) Error {
	return New().SetErrorf(msg, args...)
}

func Code(code int) Error {
	return New().SetCode(code)
}

func Err(err error) Error {
	return New().SetError(err)
}

func Status(status status.Status) Error {
	return New().SetStatus(status)
}

func BadRequest(err error) Error {
	return New().SetError(err).SetStatus(status.BadRequest)
}

func InternalError(err error) Error {
	return New().SetError(err).SetStatus(status.InternalServerError)
}

func UnAuthorized(err error) Error {
	return New().SetError(err).SetStatus(status.Unauthorized)
}

func Forbidden(err error) Error {
	return New().SetError(err).SetStatus(status.Forbidden)
}
