package main

import (
	"github.com/246859/ginx"
	"github.com/246859/ginx/constant/status"
	"github.com/246859/ginx/pkg/resp"
	"github.com/gin-gonic/gin"
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
		resp.Ok(ctx).Status(status.OK).Msg(status.OK.String()).Data("hello world!").Render()
	})
	err := server.Spin()
	if err != nil {
		log.Fatal(err)
	}
}
