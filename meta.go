package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// M is a group of V
type M []V

func (m M) build() MetaData {
	if len(m) == 0 {
		return emptyMetaData
	}
	metadata := newMetaData()
	for _, value := range m {
		metadata.set(value.Key, value.Val)
	}
	return metadata
}

// V is basic item in metadata
type V struct {
	Key string
	Val any
}

func (v V) Value() any {
	return v.Val
}

func (v V) Bool() bool {
	if b, ok := v.Val.(bool); ok {
		return b
	}
	return false
}

func (v V) String() string {
	if s, ok := v.Val.(string); ok {
		return s
	}
	return ""
}

func (v V) Uint() uint {
	if i, ok := v.Val.(uint); ok {
		return i
	}
	return 0
}

func (v V) Uint8() uint8 {
	if i, ok := v.Val.(uint8); ok {
		return i
	}
	return 0
}

func (v V) Uint16() uint16 {
	if i, ok := v.Val.(uint16); ok {
		return i
	}
	return 0
}

func (v V) Uint32() uint32 {
	if i, ok := v.Val.(uint32); ok {
		return i
	}
	return 0
}

func (v V) Uint64() uint64 {
	if i, ok := v.Val.(uint64); ok {
		return i
	}
	return 0
}

func (v V) Int() int {
	if i, ok := v.Val.(int); ok {
		return i
	}
	return 0
}

func (v V) Int8() int8 {
	if i, ok := v.Val.(int8); ok {
		return i
	}
	return 0
}

func (v V) Int16() int16 {
	if i, ok := v.Val.(int16); ok {
		return i
	}
	return 0
}

func (v V) Int32() int32 {
	if i, ok := v.Val.(int32); ok {
		return i
	}
	return 0
}

func (v V) Int64() int64 {
	if i, ok := v.Val.(int64); ok {
		return i
	}
	return 0
}

func (v V) Float32() float32 {
	if f, ok := v.Val.(float32); ok {
		return f
	}
	return 0
}

func (v V) Float64() float64 {
	if f, ok := v.Val.(float64); ok {
		return f
	}
	return 0
}

func (v V) Time() time.Time {
	if t, ok := v.Val.(time.Time); ok {
		return t
	}
	return time.Time{}
}

func (v V) Duration() time.Duration {
	if d, ok := v.Val.(time.Duration); ok {
		return d
	}
	return 0
}

func newMetaData() MetaData {
	return MetaData{m: make(map[string]any, 4)}
}

// MetaData is a read map store in memory
type MetaData struct {
	m map[string]any
}

func (m MetaData) Get(key string) (V, bool) {
	v, e := m.m[key]
	if !e {
		return V{}, false
	}
	return V{Key: key, Val: v}, true
}

func (m MetaData) set(k string, v any) {
	m.m[k] = v
}

func (m MetaData) applyM(meta MetaData) {
	for k, v := range meta.m {
		if !m.Has(k) {
			m.set(k, v)
		}
	}
}

func (m MetaData) ShouldGet(key string) V {
	val, _ := m.Get(key)
	return val
}

func (m MetaData) MustGet(key string) V {
	val, e := m.Get(key)
	if !e {
		panic(fmt.Sprintf("not found in metdata: %s", key))
	}
	return val
}

func (m MetaData) Has(key string) bool {
	_, has := m.Get(key)
	return has
}

func (m MetaData) Contains(v V) bool {
	get, b := m.Get(v.Key)
	if !b {
		return false
	}
	return get == v
}

func (m MetaData) String() string {
	var buf strings.Builder
	buf.WriteString("{")
	i := 0
	for k, v := range m.m {
		buf.WriteString(fmt.Sprintf("%s:%v", k, v))
		if i < len(m.m)-1 {
			buf.WriteString(",")
		}
		i++
	}
	buf.WriteString("}")
	return buf.String()
}

var emptyRouteMeta = routeMeta{MetaData: emptyMetaData}

type routeMeta struct {
	MetaData MetaData
	FullPath string
	Method   string
	// group
	Group *RouterGroup
}

const _MetaKey = "github.com/246859/ginx.metadata"

// metaDataHandler get metadata for each route from the global metadata, then store in the context
func metaDataHandler(metadata *FrozenMap[string, routeMeta]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := routeKey(ctx.Request.Method, ctx.FullPath())
		src, e := metadata.Get(key)
		if e {
			ctx.Set(_MetaKey+key, src.MetaData)
		}
	}
}

var emptyMetaData = MetaData{m: map[string]any{}}

// MetaFromCtx get metadata of route itself from context
func MetaFromCtx(ctx *gin.Context) MetaData {
	routeKey := routeKey(ctx.Request.Method, ctx.FullPath())
	metadata, exists := ctx.Get(_MetaKey + routeKey)
	if !exists {
		return emptyMetaData
	}
	return metadata.(MetaData)
}

func routeKey(method, path string) string {
	return method + ":" + path
}
