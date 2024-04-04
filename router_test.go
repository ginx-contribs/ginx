package ginx

import (
	"fmt"
	"testing"
)

func TestRouter(t *testing.T) {
	server := Default()
	root := server.RouterGroup()
	root.MGET("login", M{{"role", "guest"}, {"limit", 5}})
	user := root.MGroup("user", nil)
	user.MGET("info", M{{"role", "user"}}, nil)

	root.Walk(func(info RouteInfo) {
		t.Log(fmt.Sprintf("%+v", info))
	})
}
