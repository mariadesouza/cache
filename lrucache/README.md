#LRUCache

Package implements a LRU cache.

## Exported struct

LRUCache

## Exported Functions

###  New()

Creates a new LRUCache

  Returns:
    \*LRUCache

### (\*LRUCache) Add

Adds a key, value element to the cache

    Input:
        key string
        value string

### (\*LRUCache) Get

    Input:
      key string

    Returns:
        string - contains value corresponding to key
        bool - is true if the key exists in the cache


#Contributors
* [Maria DeSouza](maria.g.desouza@gmail.com)
