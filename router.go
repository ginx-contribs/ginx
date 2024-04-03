package ginx

import (
	"fmt"
	"github.com/246859/ginx/constant/methods"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
)

var allowMethods = []string{
	methods.Get,
	methods.Post,
	methods.Put,
	methods.Delete,
	methods.Options,
	methods.Head,
}

// RouterHandler represents a single route handler
type RouterHandler struct {
	group    *RouterGroup
	chain    gin.HandlersChain
	Method   string
	FullPath string
}

func (handler *RouterHandler) applyMeta(meta M) {
	key := routeKey(handler.Method, handler.FullPath)
	handler.group.s.metadata[key] = meta.build()
}

func (handler *RouterHandler) getMeta() MetaData {
	key := routeKey(handler.Method, handler.FullPath)
	return handler.group.s.metadata[key]
}

// RouterGroup returns metadata route group
func (s *Server) RouterGroup() *RouterGroup {
	return &RouterGroup{current: &s.engine.RouterGroup, s: s}
}

type RouterGroup struct {
	s     *Server
	group *RouterGroup

	current   *gin.RouterGroup
	handlers  []*RouterHandler
	subGroups []*RouterGroup
}

// Use same as *gin.RouterGroup.Use()
func (group *RouterGroup) Use(handlers ...gin.HandlerFunc) {
	group.current.Use(handlers...)
}

// register registers meta info into *Server.metadata
func (group *RouterGroup) applyMeta(meta M) {
	routeKey := routeKey("group", group.current.BasePath())
	group.s.metadata[routeKey] = meta.build()
}

func (group *RouterGroup) getMeta() MetaData {
	routeKey := routeKey("group", group.current.BasePath())
	return group.s.metadata[routeKey]
}

func (group *RouterGroup) Group(path string, meta M, handlers ...gin.HandlerFunc) *RouterGroup {
	// register route
	newGroup := group.current.Group(path, handlers...)
	// register metadata
	subGroup := &RouterGroup{
		group:   group,
		current: newGroup,
		s:       group.s,
	}
	subGroup.applyMeta(meta)

	group.subGroups = append(group.subGroups, subGroup)
	return subGroup
}

func (group *RouterGroup) Handle(method string, path string, meta M, handlers ...gin.HandlerFunc) *RouterHandler {
	// register route
	group.current.Handle(method, path, handlers...)
	handler := &RouterHandler{
		group:    group,
		chain:    handlers,
		Method:   method,
		FullPath: joinPaths(group.current.BasePath(), path),
	}
	// register metadata
	handler.applyMeta(meta)
	group.handlers = append(group.handlers, handler)
	return handler
}

func (group *RouterGroup) Match(methods []string, path string, meta M, handlers ...gin.HandlerFunc) []*RouterHandler {
	var hs []*RouterHandler
	for _, method := range methods {
		if !slices.Contains(allowMethods, method) {
			panic(fmt.Sprintf("not allowed method: %s", method))
		}
		hs = append(hs, group.Handle(method, path, meta, handlers...))
	}
	return hs
}

func (group *RouterGroup) GET(path string, meta M, handlers ...gin.HandlerFunc) *RouterHandler {
	return group.Handle(http.MethodGet, path, meta, handlers...)
}

func (group *RouterGroup) POST(path string, meta M, handlers ...gin.HandlerFunc) *RouterHandler {
	return group.Handle(http.MethodPost, path, meta, handlers...)
}

func (group *RouterGroup) DELETE(path string, meta M, handlers ...gin.HandlerFunc) *RouterHandler {
	return group.Handle(http.MethodDelete, path, meta, handlers...)
}

func (group *RouterGroup) PUT(path string, meta M, handlers ...gin.HandlerFunc) *RouterHandler {
	return group.Handle(http.MethodPut, path, meta, handlers...)
}

func (group *RouterGroup) OPTIONS(path string, meta M, handlers ...gin.HandlerFunc) *RouterHandler {
	return group.Handle(http.MethodOptions, path, meta, handlers...)
}

func (group *RouterGroup) HEAD(path string, meta M, handlers ...gin.HandlerFunc) *RouterHandler {
	return group.Handle(http.MethodHead, path, meta, handlers...)
}

func (group *RouterGroup) Any(path string, meta M, handlers ...gin.HandlerFunc) []*RouterHandler {
	return group.Match(allowMethods, path, meta, handlers...)
}

type RouteInfo struct {
	IsGroup bool
	Group   *RouteInfo

	Method   string
	FullPath string
	Handler  gin.HandlerFunc
	Meta     MetaData
}

// Walk group and handlers info, include subgroup
func (group *RouterGroup) Walk(walkFn func(info RouteInfo)) {
	infoList := make([]RouteInfo, 0, len(group.handlers)+1)

	// append subGroup info
	groupInfo := RouteInfo{
		IsGroup:  true,
		FullPath: group.current.BasePath(),
		Meta:     group.getMeta(),
	}
	infoList = append(infoList, groupInfo)

	// append route info
	for _, handler := range group.handlers {
		infoList = append(infoList, RouteInfo{
			IsGroup:  false,
			Method:   handler.Method,
			FullPath: handler.FullPath,
			Handler:  lastHandler(handler.chain),
			Meta:     handler.getMeta(),
			Group:    &groupInfo,
		})
	}

	// walk sub handlers info
	for _, info := range infoList {
		walkFn(info)
	}

	// then walk subgroups recursively
	for _, subGroup := range group.subGroups {
		subGroup.Walk(walkFn)
	}
}
