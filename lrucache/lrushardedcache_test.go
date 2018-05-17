package lrucache

import (
	"testing"
)

// Uses the same testTable as the lruCache
func TestLRUShardedCacheSetGet(t *testing.T) {
	lruShardedCache := NewShardedCache(1, 120)
	for _, test := range testTable {
		lruShardedCache.Set(test.keySet, test.value)
		value, ok := lruShardedCache.Get(test.keyGet)
		if ok != test.expectedResult {
			t.Errorf("LRUShardedCache returned = %v; want: %v", ok, test.expectedResult)
		} else if value != test.value {
			t.Errorf("LRUShardedCache returned = %v; want: %v", value, test.value)
		}
	}
}

func TestLRUShardedCacheRemove(t *testing.T) {
	lruShardedCache := NewShardedCache(2, 120)
	lruShardedCache.Set("pincode", 1234)
	ok := lruShardedCache.Remove("pincode")
	expectedResult := true
	if !ok {
		t.Errorf("LRUShardedCache.Remove returned = %v; want: %v", ok, expectedResult)
	}
	_, ok = lruShardedCache.Get("pincode")
	if ok == true {
		t.Errorf("LRUShardedCache.Get returned true -  LRUCache.Remove failed to a removed entry")
	}
}
