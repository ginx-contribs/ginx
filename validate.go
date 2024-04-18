package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"github.com/go-playground/validator/v10"
)

// ValidateHandler will be called if validate failed.
type ValidateHandler func(ctx *gin.Context, val any, err error)

var DefaultValidateHandler ValidateHandler = func(ctx *gin.Context, val any, err error) {
	response := resp.Fail(ctx).Error(err)
	if _, ok := err.(validator.ValidationErrors); ok {
		response = response.Msg("invalid parameters")
	} else {
		response = response.Status(status.InternalServerError).Msg(status.InternalServerError.String())
	}
	response.JSON()
}

func ShouldValidateWith(ctx *gin.Context, val any, binding binding.Binding) error {
	err := ctx.ShouldBindWith(val, binding)
	if err != nil {
		DefaultValidateHandler(ctx, val, err)
		return err
	}
	return nil
}

func ShouldValidate(ctx *gin.Context, val any) error {
	err := ctx.ShouldBind(val)
	if err != nil {
		DefaultValidateHandler(ctx, val, err)
		return err
	}
	return nil
}

func ShouldValidateJSON(ctx *gin.Context, val any) error {
	err := ctx.ShouldBind(val)
	if err != nil {
		DefaultValidateHandler(ctx, val, err)
		return err
	}
	return nil
}

func ShouldValidateQuery(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindQuery(val)
	if err != nil {
		DefaultValidateHandler(ctx, val, err)
		return err
	}
	return nil
}

func ShouldValidateURI(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindUri(val)
	if err != nil {
		DefaultValidateHandler(ctx, val, err)
		return err
	}
	return nil
}

func ShouldValidateHeader(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindHeader(val)
	if err != nil {
		DefaultValidateHandler(ctx, val, err)
		return err
	}
	return nil
}

func ShouldValidateXML(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindXML(val)
	if err != nil {
		DefaultValidateHandler(ctx, val, err)
		return err
	}
	return nil
}

func ShouldValidateYAML(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindYAML(val)
	if err != nil {
		DefaultValidateHandler(ctx, val, err)
		return err
	}
	return nil
}

func ShouldValidateTOML(ctx *gin.Context, val any) error {
	err := ctx.ShouldBindTOML(val)
	if err != nil {
		DefaultValidateHandler(ctx, val, err)
		return err
	}
	return nil
}
