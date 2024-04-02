# ginx
ginx is a simple gin enhancer, which has those features:

* lightweight and more convenient
* graceful shutdown
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