#redisproxy

Package that implements the redisproxy.

## Exported struct

Redisproxy

## Exported Functions

###  New

    Input:
      redisServer string
      port string
      cacheCapacity int
      cacheExpirySeconds int

  Returns:


### (\*Redisproxy) Close

closes the Redis connection

### Get
    Input:
      key string

    Returns:
      string containing value

#Contributors
* [Maria DeSouza](maria.g.desouza@gmail.com)
