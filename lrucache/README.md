#LRUCache

Package implements a LRU cache.

## Exported struct

LRUCache

## Exported Functions

### LRUCache

This is the main LRU cache implementation

####  New()

Creates a new LRUCache

  Returns:
    \*LRUCache

#### (\*LRUCache) Add

Adds a key, value element to the cache

    Input:
        key string
        value string

#### (\*LRUCache) Get

    Input:
      key string

    Returns:
        string - contains value corresponding to key
        bool - is true if the key exists in the cache

#### (\*LRUCache) Remove

Input:
  key string

### LRUShardedCache

This is the Sharded LRU cache implementation. It creates multiple shards of the LRUCache based on the hash of the key. 

#Contributors
* [Maria DeSouza](maria.g.desouza@gmail.com)
