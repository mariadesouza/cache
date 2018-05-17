package lrucache

import (
	"crypto/sha1"
	"fmt"
	"time"
)

const (
	numOfShards = 256
)

// This is a sharded version of the cache.
/*

A call to set in the lrucache would block calls to Get.
I have added sharding to reduce write locks. The keys are distributed over 256 shards.
This will improve the time for retrieval and will reduce blocking especially when concurrency is in play.

There are 256 shards, and each shard capacity will be as set when creating the new instance of the cache.
The shards keys are pre-allocated since we know all 256 combinations.
Although, this will increase memory usage the performance time will improve.

*/

//LRUShardedCache : map of LRUcache shards. map key is hash of the cache item key
type LRUShardedCache map[string]*LRUCache

// NewShardedCache : create new LRUCache
//capacity is multiplied by 256 since there are 256 shards
func NewShardedCache(capacity int, ttlSeconds time.Duration) *LRUShardedCache {
	if capacity == 0 { // set minimum capacity to 1
		capacity = 1
	}
	if ttlSeconds == 0 {
		ttlSeconds = 180
	}

	c := make(LRUShardedCache, numOfShards)
	for i := 0; i < numOfShards; i++ {
		c[fmt.Sprintf("%02x", i)] = New(capacity, ttlSeconds)
	}
	return &c
}

//Set : add node to LRUCache
func (c *LRUShardedCache) Set(key string, value interface{}) {
	shard := c.getShard(key)
	if shard == nil {
		return
	}
	shard.Set(key, value)
	//fmt.Println(shard)
}

// getShard : getShard from key hash
func (c *LRUShardedCache) getShard(key string) *LRUCache {
	hasher := sha1.New()
	hasher.Write([]byte(key))
	shardKey := fmt.Sprintf("%x", hasher.Sum(nil))[0:2]
	//fmt.Println(key, shardKey)
	return (*c)[shardKey]
}

// Get : Get value from key
func (c *LRUShardedCache) Get(key string) (interface{}, bool) {
	shard := c.getShard(key)
	if shard != nil {
		return shard.Get(key)
	}
	return "", false
}

//Remove : Remove element with key
func (c *LRUShardedCache) Remove(key string) bool {
	shard := c.getShard(key)
	if shard == nil {
		return false
	}
	return shard.Remove(key)
}
