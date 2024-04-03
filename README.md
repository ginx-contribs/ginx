# ginx
ginx is a simple gin enhancer, which has those features:

* lightweight and more convenient
* graceful shutdown
* walk routes and store metadata for each route
* hooks at `BeforeStarting`, `AfterSarted`, `OnShutdown`
* integrated with many useful middleware, like `ratelimit`,`cache`, `accesslog`, `recovery` and so on.

## install
```bash
go get github.com/246859/ginx@latest
```

## usage
see more examples in [ginx examples](https://github.com/246859/ginx/tree/main/examples/).

### quick start
```go
package main

import (
	"github.com/246859/ginx"
	"log"
)

func main() {
	server := ginx.New()
	err := server.Spin()
	if err != nil {
		log.Fatal(err)
	}
}

```

### meta group
use the meta group and walk with func.
```go
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
```

### use middleware
see more middlewares at [ginx middlewares](https://github.com/246859/ginx/tree/main/contribs/).
```go
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

```