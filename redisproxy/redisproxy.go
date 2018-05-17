package redisproxy

import (
	"time"

	"github.com/mariadesouza/redisproxyserver/lrucache"
)

//RedisProxy :
type RedisProxy struct {
	redisConn interface{}
	cache     *lrucache.LRUCache
	//cache *lrucache.LRUShardedCache
}

//New : create new redisproxy object
func New(redisServer string, port string, cacheCapacity int, cacheExpirySeconds int64) (*RedisProxy, error) {
	var redisproxy RedisProxy
	var err error
	redisproxy.redisConn, err = newRedisConnection(redisServer, port)
	if err != nil {
		return nil, err
	}
	expiryTime := time.Duration(cacheExpirySeconds)
	redisproxy.cache = lrucache.New(cacheCapacity, expiryTime)
	//redisproxy.cache = lrucache.NewShardedCache(cacheCapacity, expiryTime)
	return &redisproxy, nil
}

//Close : Close connection to Redis
func (r *RedisProxy) Close() {
	r.redisConn.(*redisServerConn).CloseConnection()
}

// Get : fetch value of key from cache/redis
func (r *RedisProxy) Get(key string) (interface{}, error) {

	// look into LRUCache
	if value, ok := r.cache.Get(key); ok {
		return value, nil
	}
	//does not exist - look in redis
	value, err := r.getRedisValue(key)
	if err != nil {
		return "", err
	}

	// found it - set it in cache
	r.cache.Set(key, value)
	return value, nil
}

func (r *RedisProxy) getRedisValue(key string) (string, error) {
	redisConn, ok := r.redisConn.(*redisServerConn)
	if ok {
		err := redisConn.Send("GET", key)
		if err != nil {
			return "", err
		}
		response, err := r.redisConn.(*redisServerConn).Receive()
		if err != nil {
			return "", err
		}
		value := string(response[:])
		return value, nil
	}
	return "", nil
}
