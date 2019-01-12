package redisproxy

import (
	"testing"
	"time"

	"github.com/mariadesouza/cache/lrucache"
)

// This will mock the RedisProxy struct from the redisproxy package
// This is so our unit test doesnt have to make an actual connection to a redis instance
type redisServerConnector interface {
	Send(args ...string) error
}

type mockRedis struct {
}

var _ redisServerConnector = (*mockRedis)(nil)

func (f *mockRedis) Send(args ...string) error {
	return nil
}

func setupNewTest(cacheCapacity int, cacheExpirySeconds int64) (*RedisProxy, string, string) {
	var redisproxy RedisProxy
	expiryTime := time.Duration(cacheExpirySeconds)
	//redisproxy.cache = lrucache.New(cacheCapacity, expiryTime)
	redisproxy.cache = lrucache.NewShardedCache(cacheCapacity, expiryTime)
	redisproxy.redisConn = &mockRedis{}
	value := "apple"
	key := "fruit"
	redisproxy.cache.Set(key, value)
	return &redisproxy, key, value
}

func TestGetSuccess(t *testing.T) {
	redisproxy, key, value := setupNewTest(5, 10)
	redisproxy.cache.Set(key, value)
	res, err := redisproxy.Get(key)
	if err != nil {
		t.Errorf("redisproxy returned error %v", err)
	} else if value != res {
		t.Errorf("redisproxy returned = %v; want: %v", res, value)
	}
}

func TestGetNonExistingKey(t *testing.T) {
	redisproxy, _, _ := setupNewTest(5, 10)
	res, _ := redisproxy.Get("notthere")
	if res != "" {
		t.Errorf("redisproxy returned = %v; want: ", res)
	}
}
