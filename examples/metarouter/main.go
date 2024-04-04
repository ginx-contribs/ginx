package main

import (
	"fmt"
	"github.com/246859/ginx"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
)

func main() {
	server := ginx.Default()
	root := server.RouterGroup()
	root.MGET("login", ginx.M{{"role", "guest"}, {"limit", 5}})
	user := root.MGroup("user", nil)
	user.MGET("info", ginx.M{{"role", "user"}}, func(ctx *gin.Context) {
		// get metadata from context
		metaData := ginx.MetaFromCtx(ctx)
		slog.Info(metaData.ShouldGet("role").String())
	})

	// walk root router
	root.Walk(func(info ginx.RouteInfo) {
		slog.Info(fmt.Sprintf("%s %s", info.FullPath, info.Meta))
	})

	err := server.Spin()
	if err != nil {
		log.Fatal(err)
	}
}
