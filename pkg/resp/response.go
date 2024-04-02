package resp

import (
	"errors"
	"github.com/246859/ginx/constant/status"
	"github.com/246859/ginx/pkg/resp/errs"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response represents a http json response
type Response struct {
	body struct {
		Code int    `json:"code,omitempty"`
		Data any    `json:"data,omitempty"`
		Msg  string `json:"msg,omitempty"`
	}

	status status.Status
	err    error

	ctx *gin.Context
}

func (resp *Response) Code(code int) *Response {
	resp.body.Code = code
	return resp
}

func (resp *Response) Data(data any) *Response {
	resp.body.Data = data
	return resp
}

func (resp *Response) Msg(msg string) *Response {
	resp.body.Msg = msg
	return resp
}

func (resp *Response) Error(err error) *Response {
	resp.err = err
	return resp
}

func (resp *Response) Status(status status.Status) *Response {
	resp.status = status
	return resp
}

func (resp *Response) Render() {
	ctx := resp.ctx
	if ctx == nil {
		return
	}

	if resp.body.Code == 0 {
		resp.body.Code = resp.status.Code()
	}

	if resp.err != nil {
		// if is status error
		var statusErr errs.Error
		if ok := errors.As(resp.err, &statusErr); ok {
			resp.status = statusErr.Status
		}

		// overlay the body msg
		if resp.body.Msg == "" {
			if statusErr.Status != http.StatusInternalServerError {
				resp.body.Msg = statusErr.Error()
			} else {
				resp.body.Msg = status.InternalServerError.String()
			}
		}
		resp.ctx.Error(resp.err)
	}

	// fall back message
	if resp.status.Code() >= 300 && resp.body.Msg == "" {
		resp.body.Msg = resp.status.String()
	}

	ctx.JSON(resp.status.Code(), resp.body)
}

func New(ctx *gin.Context) *Response {
	return &Response{ctx: ctx}
}

// Ok response with status code 200
func Ok(ctx *gin.Context) *Response {
	return &Response{ctx: ctx, status: status.OK}
}

// Fail response with status code 400
func Fail(ctx *gin.Context) *Response {
	return &Response{ctx: ctx, status: status.BadRequest}
}

// Forbidden response with status code 403
func Forbidden(ctx *gin.Context) *Response {
	return &Response{ctx: ctx, status: status.Forbidden}
}

// UnAuthorized response with status code 401
func UnAuthorized(ctx *gin.Context) *Response {
	return &Response{ctx: ctx, status: status.Forbidden}
}
