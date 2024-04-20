package ginx

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"sync/atomic"
)

// FrozenMap is a thread-safe map, write operations will cause panic after map is frozen.
type FrozenMap[K comparable, V any] struct {
	readOnly atomic.Bool
	m        cmap.ConcurrentMap[K, V]
}

// Frozen make map be immutable
func (r *FrozenMap[K, V]) Frozen() {
	if r.readOnly.Load() == false {
		r.readOnly.Store(true)
	}
}

func (r *FrozenMap[K, V]) Set(k K, v V) {
	if r.readOnly.Load() {
		panic("map is frozen, write operations is not permitted")
	}
	r.m.Set(k, v)
}

func (r *FrozenMap[K, V]) Get(k K) (V, bool) {
	return r.m.Get(k)
}

func (r *FrozenMap[K, V]) Del(k K) {
	r.m.Remove(k)
}

func (r *FrozenMap[K, V]) Range(fn func(K, V)) {
	r.m.IterCb(fn)
}
