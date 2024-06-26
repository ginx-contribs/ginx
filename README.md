# ginx
![Static Badge](https://img.shields.io/badge/go-1.22-blue)
![GitHub License](https://img.shields.io/github/license/ginx-contribs/ginx)
![go reports](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)


ginx is a simple gin enhancer, features as follows:

* lightweight and more convenient
* graceful shutdown
* support walk router and store metadata
* hooks at `BeforeStarting`, `AfterSarted`, `OnShutdown`
* integrated with many useful middleware, like `ratelimit`,`recovery`, `accesslog`, `rquestId`, `cors` and so on.

## install
```bash
go get github.com/ginx-contribs/ginx@latest
```

## usage
see more examples in [ginx examples](https://github.com/246859/ginx/tree/main/examples/).

### quick start
```go
package main

import (
	"github.com/ginx-contribs/ginx"
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
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
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
```
output
```
2024/04/04 12:49:16 INFO / {}
2024/04/04 12:49:16 INFO /login {role:guest,limit:5}
2024/04/04 12:49:16 INFO /user {}
2024/04/04 12:49:16 INFO /user/info {role:user}
2024/04/04 12:49:16 INFO [GinX] server is listening on :8080
2024/04/04 12:49:27 INFO user
2024/04/04 12:49:27 INFO [GinX] status=200 method=GET cost=201.60µs ip=127.0.0.1 url=/user/info path=/user/info route=/user/info request-size=720B response-size=0B
2024/04/04 12:51:16 INFO [GinX] received stop signal, it will shutdown in 5s at latest
2024/04/04 12:51:16 INFO [GinX] server shutdown
```

### use middleware
see more middlewares at [ginx middlewares](https://github.com/246859/ginx/tree/main/contribs/).
```go
package main

import (
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx/middleware"
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
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
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
		resp.Fail(ctx).Error(statuserr.New().SetErrorf("invalid access").SetCode(1018)).JSON()
	})
	err := server.Spin()
	if err != nil {
		log.Fatal(err)
	}
}
```
