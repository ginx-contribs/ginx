package counter

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func testCounter(t *testing.T, counter Counter, key string, limit int, window time.Duration) {
	for i := 0; i < 20; i++ {
		time.Sleep(time.Millisecond * 500)
		count, err := counter.Count(context.Background(), key, limit, window)
		t.Log(i, count, err)
		assert.Nil(t, err)
		if i > limit {
			assert.EqualValues(t, count, limit)
		}
	}
}

func TestCacheCounter(t *testing.T) {
	counter := Cache()
	limit := 10
	window := time.Second * 10
	key := "a"
	testCounter(t, counter, key, limit, window)
}

func TestRedisCounter(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "192.168.48.144:6379", Password: "123456"})
	counter := Redis(client)
	limit := 10
	window := time.Second * 10
	key := "a"
	testCounter(t, counter, key, limit, window)
}
