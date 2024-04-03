package main

import (
	"fmt"
	"github.com/246859/ginx"
	"log/slog"
)

func main() {
	server := ginx.Default()
	root := server.RouterGroup()
	root.GET("login", ginx.M{{"role", "guest"}, {"limit", 5}})
	user := root.Group("user", nil)
	user.GET("info", ginx.M{{"role", "user"}}, nil)

	root.Walk(func(info ginx.RouteInfo) {
		slog.Info(fmt.Sprintf("%+v", info))
	})
}
