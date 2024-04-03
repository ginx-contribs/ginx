package ginx

import (
	"fmt"
	"testing"
)

func TestRouter(t *testing.T) {
	server := Default()
	root := server.RouterGroup()
	root.GET("login", M{{"role", "guest"}, {"limit", 5}})
	user := root.Group("user", nil)
	user.GET("info", M{{"role", "user"}}, nil)

	root.Walk(func(info RouteInfo) {
		t.Log(fmt.Sprintf("%+v", info))
	})
}
