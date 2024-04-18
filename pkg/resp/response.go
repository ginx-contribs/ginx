package resp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
)

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

// Response represents a http json response
type Response struct {
	body struct {
		Code  int    `json:"code,omitempty"`
		Data  any    `json:"data,omitempty"`
		Msg   string `json:"msg,omitempty"`
		Error string `json:"error,omitempty"`
	}

	// decide whether to show internal server error message in response body
	transparent bool

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

func (resp *Response) Transparent() *Response {
	resp.transparent = true
	return resp
}

// render the response body
func (resp *Response) render() {
	ctx := resp.ctx
	if ctx == nil {
		panic("nil *gin.Context in response")
	}

	if resp.err != nil {
		// if is status error
		var statusErr statuserr.Error
		if ok := errors.As(resp.err, &statusErr); ok {
			resp.status = statusErr.Status
			resp.body.Code = statusErr.Code
		}

		if resp.body.Error == "" {
			if resp.transparent || resp.status != status.InternalServerError {
				resp.body.Error = resp.err.Error()
			} else {
				// do not expose error msg for internal error,
				// it will be passed to the context, and will be processed by others
				resp.body.Error = status.InternalServerError.String()
			}
		}

		// append error into context
		resp.ctx.Error(resp.err)
	}

	// code fallback
	if resp.body.Code == 0 {
		resp.body.Code = resp.status.Code()
	}

	// fall back error msg
	if resp.status.Code() >= 400 && resp.body.Error == "" {
		resp.body.Error = resp.status.String()
	}
}

func (resp *Response) JSON() {
	resp.render()
	resp.ctx.JSON(resp.status.Code(), resp.body)
}

func (resp *Response) XML() {
	resp.render()
	resp.ctx.XML(resp.status.Code(), resp.body)
}

func (resp *Response) YAML() {
	resp.render()
	resp.ctx.YAML(resp.status.Code(), resp.body)
}

func (resp *Response) TOML() {
	resp.render()
	resp.ctx.TOML(resp.status.Code(), resp.body)
}

func (resp *Response) ProtoBuf() {
	resp.render()
	resp.ctx.ProtoBuf(resp.status.Code(), resp.body)
}
