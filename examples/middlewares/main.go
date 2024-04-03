package main

import (
	"github.com/246859/ginx"
	"github.com/246859/ginx/middleware"
	"log"
	"log/slog"
	"time"
)

func main() {
	server := ginx.New(
		ginx.WithNoRoute(middleware.NoRoute()),
		ginx.WithNoMethod(middleware.NoMethod()),
		ginx.WithMiddlewares(
			middleware.Logger(slog.Default(), "ginx"),
			middleware.RateLimit(nil, nil),
			middleware.CacheMemory("cache", time.Second),
		),
	)

	err := server.Spin()
	if err != nil {
		log.Fatal(err)
	}
}
