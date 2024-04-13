package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"github.com/ginx-contribs/ginx/pkg/resp/errs"
	"log"
)

func main() {
	server := ginx.Default()
	root := server.RouterGroup()
	// {"code":200,"data":"hello world!","msg":"ok"}
	root.GET("/hello", func(ctx *gin.Context) {
		resp.Ok(ctx).Data("hello world!").Msg("ok").JSON()
	})

	// {"code":1018,"error":"invalid access"}
	root.GET("/error", func(ctx *gin.Context) {
		resp.Fail(ctx).Error(errs.CodeWrapError(1018, "invalid access")).JSON()
	})
	err := server.Spin()
	if err != nil {
		log.Fatal(err)
	}
}
