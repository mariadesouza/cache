package lrucache

import (
	"testing"
	"time"
)

// this test table is shared by both lrucache and lrushardedcache_test
var testTable = []struct {
	name           string
	keySet         string
	keyGet         string
	value          interface{}
	expectedResult bool
}{
	{"stringGetSuccess", "fruit", "fruit", "apple", true},
	{"stringGetSuccess", "counters", "counters", 10, true},
	{"stringGetSuccess", "friends", "friends", "jane", true},
	{"stringGetSuccess", "test2", "test2", "hey", true},
	{"stringGetSuccess", "key1", "key1", "hello", true},
	{"stringGetFail", "vegetable", "notthere", "", false},
}

func TestLRUCacheSetGet(t *testing.T) {
	lruCache := New(len(testTable), 120)
	for _, test := range testTable {
		lruCache.Set(test.keySet, test.value)
		value, ok := lruCache.Get(test.keyGet)
		if ok != test.expectedResult {
			t.Errorf("LRUCache.Get returned = %v; want: %v", ok, test.expectedResult)
		} else if value != test.value {
			t.Errorf("LRUCache.Get value returned = %v; want: %v", value, test.value)
		}
	}
}

func TestLRUCacheEviction(t *testing.T) {
	lruCache := New(len(testTable)-1, 120)
	//add all
	for _, test := range testTable {
		lruCache.Set(test.keySet, test.value)
	}
	//check for first one (Least Recently Used)
	expectedResult := false
	_, ok := lruCache.Get("fruit")
	if ok != expectedResult { // should have booted off the first key as the cacpacity is one less
		t.Errorf("LRUCache.Get LRU evicted test returned = %v; want: %v", ok, expectedResult)
	}
}

func TestLRUCachePromotion(t *testing.T) {
	lruCache := New(len(testTable)-1, 120)
	//add all
	for n, test := range testTable {
		lruCache.Set(test.keySet, test.value)
		if n == int(len(testTable)/2) {
			_, _ = lruCache.Get(testTable[0].keyGet)
		}
	}
	//check for first one taht was promoted (Least Recently Used)
	expectedResult := true
	_, ok := lruCache.Get("fruit")
	if ok != expectedResult {
		t.Errorf("LRUCache.Get LRU evicted test returned = %v; want: %v", ok, expectedResult)
	}

	expectedResult2 := false
	_, ok = lruCache.Get("counters")
	if ok != expectedResult2 { // should have booted off the first key was promoted and this was LRU
		t.Errorf("LRUCache.Get LRU evicted test returned = %v; want: %v", ok, expectedResult2)
	}

}

func TestLRUCacheRemove(t *testing.T) {
	lruCache := New(1, 120)
	lruCache.Set("pincode", 1234)
	ok := lruCache.Remove("pincode")
	expectedResult := true
	if !ok {
		t.Errorf("LRUCache.Remove returned = %v; want: %v", ok, expectedResult)
	}
	_, ok = lruCache.Get("pincode")
	if ok == true {
		t.Errorf("LRUCache.Get returned true -  LRUCache.Remove failed to a removed entry")
	}
}

func TestLRUCacheExpiration(t *testing.T) {
	key := "jane"
	expectedResult := false
	lruCache := New(5, 1)
	lruCache.Set(key, "buddy")
	time.Sleep(2 * time.Second)
	_, ok := lruCache.Get(key)
	if ok != expectedResult {
		t.Errorf("LRUCache expration test returned = %v; want: %v", ok, expectedResult)
	}
}
