//package lrucache ...

package lrucache

import (
	"container/list"
	"sync"
	"time"
)

type node struct {
	key        string
	value      interface{}
	expiration int64
}

func (n *node) Expired() bool {
	if n.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > n.expiration
}

// LRUCache : least recently used cache
// doubleList is a double linked that makes  add/remove O(1)
// the cache is a hash map of keys to elements makes get O(1)
type LRUCache struct {
	capacity   int
	ttlSeconds time.Duration // in seconds
	cache      map[string]*list.Element
	doubleList *list.List
	mu         sync.RWMutex
}

// New : create new LRUCache
func New(capacity int, ttlSeconds time.Duration) *LRUCache {
	if capacity == 0 { // set minimum capacity to 1
		capacity = 1
	}
	if ttlSeconds == 0 {
		ttlSeconds = 180
	}
	return &LRUCache{
		capacity:   capacity,
		ttlSeconds: ttlSeconds,
		cache:      make(map[string]*list.Element),
		doubleList: list.New(),
	}
}

//Add : add node to LRUCache
func (c *LRUCache) Add(key string, value interface{}) {

	if c.cache == nil {
		c.cache = make(map[string]*list.Element)
		c.doubleList = list.New()
	}
	c.mu.Lock()
	// check if it exists already
	if element, ok := c.cache[key]; ok {
		//found item - promote it
		c.doubleList.MoveToFront(element)
		element.Value.(*node).key = key
		element.Value.(*node).value = value
		c.mu.Unlock()
		return
	}
	// make New - Add it to front of the list
	expireTime := time.Now().Add(time.Second * c.ttlSeconds).UnixNano()
	element := c.doubleList.PushFront(&node{key, value, expireTime})
	c.cache[key] = element
	// remove the least recently used element from the back of the list
	if c.doubleList.Len() > c.capacity {
		last := c.doubleList.Back()
		c.removeNode(last)
	}
	c.mu.Unlock()
}

//Remove : Remove a key
func (c *LRUCache) Remove(key string) bool{

	if c.cache == nil {
		return false
	}

	if element, ok := c.cache[key]; ok {
		c.removeNode(element)
		return true
	}

	return false
}

func (c *LRUCache) removeNode(e *list.Element) {
	key := e.Value.(*node).key
	c.doubleList.Remove(e)
	delete(c.cache, key)
}

// Get : fetch value for key from cache if exists
func (c *LRUCache) Get(key string) (interface{}, bool) {
	if c.cache == nil {
		return "", false
	}
	c.mu.RLock()
	if element, ok := c.cache[key]; ok {
		n := element.Value.(*node)
		if n.expiration > 0 {
				if time.Now().UnixNano() > n.expiration { // has expired remove it
				c.removeNode(element)
				c.mu.RUnlock()
				return "", false
			}
		}
		//found item - promote it
		c.doubleList.MoveToFront(element)
		c.mu.RUnlock()
		return element.Value.(*node).value, ok
	}
	c.mu.RUnlock()
	return "", false
}
