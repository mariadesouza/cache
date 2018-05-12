package redisproxy

import (
	"fmt"
	"time"

	"github.com/mariadesouza/redisproxyserver/lrucache"
)

//Redisproxy :
type Redisproxy struct {
	// Multiple clients are able to concurrently connect to the proxy
	//up to some configurable maximum limit  without adversely impacting the
	//functional behaviour of the proxy.
	//When multiple clients make concurrent requests to the proxy, it is acceptable
	//for them to be processed sequentially  i.e. a request from the second only starts processing after the first request
	//has completed and a response has been returned to the first client .
	////do a connection pool??
	redisConn *redisServerConn
	cache     *lrucache.LRUCache
	//cacheCapacity int
	//cacheExpiry   int
}

//New : create new redisproxy object
func New(redisServer string, port string, cacheCapacity int, cacheExpirySeconds int64) (*Redisproxy, error) {
	var redisproxy Redisproxy
	var err error
	redisproxy.redisConn, err = newRedisConnection(redisServer, port)
	if err != nil {
		return nil, err
	}
	expiryTime := time.Duration(cacheExpirySeconds)
	redisproxy.cache = lrucache.New(cacheCapacity, expiryTime)
	return &redisproxy, nil
}

//Close : Close connection to Redis
func (r *Redisproxy) Close() {
	r.redisConn.CloseConnection()
}

// Get : fetch value of key from cache/redis
func (r *Redisproxy) Get(key string) (string, error) {

	// look into LRUCache
	if value, ok := r.cache.Get(key); ok {
		return value, nil
	}
	fmt.Println("not found - lets look in redis")
	//does not exist - look in redis
	value, err := r.getRedisValue(key)
	if err != nil {
		return "", err
	}

	// found it - set it in cache
	r.cache.Add(key, value)
	return value, nil
}

func (r *Redisproxy) getRedisValue(key string) (string, error) {
	err := r.redisConn.Send("GET", key)
	if err != nil {
		return "", err
	}
	response, err := r.redisConn.Receive()
	if err != nil {
		return "", err
	}

	value := string(response[:])

	return value, nil
}
