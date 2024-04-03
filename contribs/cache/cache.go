package cache

import (
	gincache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewMemStore(ttl time.Duration) *persist.MemoryStore {
	return persist.NewMemoryStore(ttl)
}

// NewRedisStore create a redis memory store with redis client
func NewRedisStore(redisClient *redis.Client) *RedisStore {
	return &RedisStore{
		RedisClient: redisClient,
	}
}

type Options struct {
	CacheOpts []gincache.Option
	Prefix    string
	TTl       time.Duration
	Store     persist.CacheStore
	Strategy  gincache.GetCacheStrategyByRequest
}

type Option func(options *Options)

func WithTTL(ttl time.Duration) Option {
	return func(options *Options) {
		options.TTl = ttl
	}
}

func WithPrefix(prefix string) Option {
	return func(options *Options) {
		options.Prefix = prefix
	}
}

func WithStrategy(strategy gincache.GetCacheStrategyByRequest) Option {
	return func(options *Options) {
		options.Strategy = strategy
	}
}

func WithStore(store persist.CacheStore) Option {
	return func(options *Options) {
		options.Store = store
	}
}

func WithCacheOption(cacheOpts ...gincache.Option) Option {
	return func(options *Options) {
		options.CacheOpts = cacheOpts
	}
}

// Cache returns a cache handler
func Cache(opts ...Option) gin.HandlerFunc {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	if options.Prefix == "" {
		options.Prefix = "ginx:"
	}

	if options.TTl == 0 {
		options.TTl = time.Second * 2
	}

	if options.Strategy == nil {
		options.Strategy = CacheByUri(true)
	}

	if options.Store == nil {
		options.Store = NewMemStore(options.TTl)
	}

	options.CacheOpts = append(options.CacheOpts, gincache.WithCacheStrategyByRequest(options.Strategy), gincache.WithPrefixKey(options.Prefix))

	return gincache.Cache(options.Store, options.TTl, options.CacheOpts...)
}
