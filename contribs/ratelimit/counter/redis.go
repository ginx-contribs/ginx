package counter

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

const countLuaScript = `
local count = 0
count = redis.call('get', KEYS[1])
-- check threshold
if not count then
    redis.call('set', KEYS[1], 1)
    redis.call('expire', KEYS[1], ARGV[2])
end
-- increase count
if count and tonumber(count) < tonumber(ARGV[1]) then
    redis.call('incr', KEYS[1])
end
count = redis.call('get',KEYS[1])
return tonumber(count)
`

func Redis(client *redis.Client) *RedisCounter {
	return &RedisCounter{client: client}
}

// RedisCounter implements Counter by redis lua script atomic operations
type RedisCounter struct {
	client *redis.Client
}

func (r *RedisCounter) Count(ctx context.Context, key string, limit int, window time.Duration) (int, error) {
	result := r.client.Eval(ctx, countLuaScript, []string{key}, []any{limit, int(window.Seconds())})
	return result.Int()
}
