package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"log"
)

func main() {
	server := ginx.Default()
	root := server.RouterGroup()
	// {
	// 	 "code": 200,
	// 	 "data": "hello world!",
	// 	 "msg": "OK"
	// }
	root.GET("/hello", nil, func(ctx *gin.Context) {
		resp.Ok(ctx).Status(status.OK).Msg(status.OK.String()).Data("hello world!").JSON()
	})
	err := server.Spin()
	if err != nil {
		log.Fatal(err)
	}
}
