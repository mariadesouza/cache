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
	items      map[string]*list.Element
	doubleList *list.List
	lock       *sync.RWMutex
}

// New  : create a new LRUcache
func New(capacity int, ttlSeconds time.Duration) *LRUCache {
	if capacity == 0 { // set minimum capacity to 1
		capacity = 5
	}
	if ttlSeconds == 0 {
		ttlSeconds = 180
	}

	return &LRUCache{
		capacity:   capacity,
		ttlSeconds: ttlSeconds,
		items:      make(map[string]*list.Element, capacity),
		doubleList: list.New(),
		lock:       new(sync.RWMutex),
	}
}

//Set : Set key value in cache.
func (s *LRUCache) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// check if it exists already
	if element, ok := s.items[key]; ok {
		// update value
		element.Value.(*node).value = value
		//found item - promote it
		s.doubleList.MoveToFront(element)
		return
	}
	// remove the least recently used element from the back of the list if already at capacity
	if s.doubleList.Len() == s.capacity {
		//fmt.Println("exceeded capacity", s.doubleList.Len(), s.capacity)
		last := s.doubleList.Back()
		s.removeNode(last)
	}
	// make New - Add it to front of the list
	expireTime := time.Now().Add(time.Second * s.ttlSeconds).UnixNano()
	element := s.doubleList.PushFront(&node{key, value, expireTime})

	s.items[key] = element
}

// Get : Get value of key from cache
func (s *LRUCache) Get(key string) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.items == nil {
		return "", false
	}
	if element, ok := s.items[key]; ok {
		n := element.Value.(*node)
		if n.expiration > 0 {
			if time.Now().UnixNano() > n.expiration { // has expired remove it
				s.removeNode(element)
				return "", false
			}
		}
		//found item - promote it
		s.doubleList.MoveToFront(element)
		return element.Value.(*node).value, ok
	}
	return "", false
}

//Remove : Remove  from cache
func (s *LRUCache) Remove(key string) bool {
	if element, ok := s.items[key]; ok {
		s.removeNode(element)
		return true
	}
	return false
}

func (s *LRUCache) removeNode(e *list.Element) {
	key := e.Value.(*node).key
	s.doubleList.Remove(e)
	delete(s.items, key)
}
