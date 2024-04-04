# ginx
ginx is a simple gin enhancer, features as follows:

* lightweight and more convenient
* graceful shutdown
* support walk router and store metadata
* hooks at `BeforeStarting`, `AfterSarted`, `OnShutdown`
* integrated with many useful middleware, like `ratelimit`,`recovery`, `accesslog`, `rquestId`, `cors` and so on.

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

### response
ginx provides a unified response body and supports chain calling
```go
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
		resp.Ok(ctx).Status(status.OK).Msg(status.OK.String()).Data("hello world!").JSON()
	})
	err := server.Spin()
	if err != nil {
		log.Fatal(err)
	}
}

```
