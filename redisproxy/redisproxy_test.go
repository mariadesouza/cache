package redisproxy

import(
  "testing"
  "time"
  "github.com/mariadesouza/redisproxyserver/lrucache"
)

// This will mock the RedisProxy struct from the redisproxy package
// This is so our unit test doesnt have to make an actual connection to a redis instance
type redisServerConnector interface {
  Send(args ...string) error
}

type mockRedis struct {
}

var _ redisServerConnector = (*mockRedis)(nil)

func (f *mockRedis) Send(args ...string) error{
  return nil
}

func NewTest(cacheCapacity int, cacheExpirySeconds int64) (*RedisProxy, error) {
  var redisproxy RedisProxy
	expiryTime := time.Duration(cacheExpirySeconds)
	redisproxy.cache = lrucache.New(cacheCapacity, expiryTime)
  redisproxy.redisConn = &mockRedis{}
	return &redisproxy, nil
}

func TestGet(t *testing.T){
  value:= "apple"
  key := "fruit"
  redisproxy, _ := NewTest(5,10)
  redisproxy.cache.Add(key, value)
  res, err := redisproxy.Get(key)
  if err != nil {
    t.Errorf("cache returned error %v", err)
  } else if value != res {
    t.Errorf("cache returned = %v; want: %v", value, res)
  }

}
