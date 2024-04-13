package errs

import (
	"errors"
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

func (e Error) SetError(err error) Error {
	e.Err = err
	return e
}

func (e Error) SetStatus(status status.Status) error {
	e.Status = status
	return e
}

func (e Error) Error() string {
	return e.Err.Error()
}

func New() Error {
	return Error{}
}

func CodeWrapError(code int, errMsg string) Error {
	return Error{Code: code, Err: errors.New(errMsg)}
}

func CodeError(code int, err error) Error {
	return Error{Code: code, Err: err}
}
func StatusError(err error, status status.Status) Error {
	return Error{Err: err, Status: status}
}

func BadRequest(err error) Error {
	return StatusError(err, status.BadRequest)
}

func InternalError(err error) Error {
	return StatusError(err, status.InternalServerError)
}

func UnAuthorized(err error) Error {
	return StatusError(err, status.Unauthorized)
}

func Forbidden(err error) Error {
	return StatusError(err, status.Forbidden)
}
