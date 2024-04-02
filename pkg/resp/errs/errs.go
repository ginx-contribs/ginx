package errs

import (
	"github.com/246859/ginx/constant/status"
)

// Error represents http response error, which is along with http status code,
// it used to decide how to show error message in response.
type Error struct {
	Err    error
	Status status.Status
}

func (e Error) Error() string {
	return e.Err.Error()
}

func Wrap(err error, status status.Status) Error {
	return Error{Err: err, Status: status}
}

func BadRequest(err error) Error {
	return Wrap(err, status.BadRequest)
}

func InternalError(err error) Error {
	return Wrap(err, status.InternalServerError)
}

func UnAuthorized(err error) Error {
	return Wrap(err, status.Unauthorized)
}

func Forbidden(err error) Error {
	return Wrap(err, status.Forbidden)
}
