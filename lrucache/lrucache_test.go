package lrucache

import(
  "testing"
  "time"
)

var testTable = []struct {
  name string
  keySet string
  keyGet string
  value interface{}
  expectedResult bool
}{
  {"stringGetSuccess", "fruit", "fruit", "apple", true},
  {"stringGetFail", "vegetable", "notthere", "", false},
}

func TestLRUCacheGet(t *testing.T){
  for _, test := range testTable {
    lruCache := New(5,120)
    lruCache.Add(test.keySet, test.value)
    value, ok := lruCache.Get(test.keyGet)
    if ok != test.expectedResult {
      t.Errorf("cache returned = %v; want: %v", ok, test.expectedResult)
    } else if value != test.value {
      t.Errorf("cache returned = %v; want: %v", value, test.value)
    }
  }
}

func TestLRUCacheRemove(t *testing.T){
  lruCache := New(5,120)
	lruCache.Add("pincode", 1234)

	ok := lruCache.Remove("pincode")
	if ok != true{
		t.Fatal("TestLRUCacheRemove returned a removed entry")
	}
}

func TestLRUCacheExpiration(t *testing.T){
//  ttl := 10
  key := "jane"
  expectedResult := false

  lruCache := New(5,1)
  lruCache.Add(key, "buddy")
  time.Sleep(2*time.Second)
  _, ok := lruCache.Get(key)
	if ok != expectedResult {
		  t.Errorf("cache returned = %v; want: %v", ok, expectedResult)
	}
}
