package main

import (
	"fmt"
	"github.com/246859/ginx"
	"log/slog"
)

func main() {
	server := ginx.Default()
	root := server.RouterGroup()
	root.MGET("login", ginx.M{{"role", "guest"}, {"limit", 5}})
	user := root.MGroup("user", nil)
	user.MGET("info", ginx.M{{"role", "user"}}, nil)

	root.Walk(func(info ginx.RouteInfo) {
		slog.Info(fmt.Sprintf("%s %s", info.FullPath, info.Meta))
	})
}
