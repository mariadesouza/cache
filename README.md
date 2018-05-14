# Redis Proxy webserver

A simple redis proxy server that uses an in memory LRU cache to speed up data access from redis. The cache items also have an expiration to keep data up to date.

# Design

## Server

The redisproxy server handles the GET request for a key. There are two packages implemented that are used by the server to implement the caching and redis connection management.

## lrucache

The LRU cache is implemented using a doubly linked list and a hashmap. This makes the algorithmic complexity for retrieval from the hashmap O(1) and the add/delete operations O(1). This allows the It is guarded by a reader/writer mutual exclusion lock that will help prevent contention when high read rates occur concurrently. Each item has an expiration. If an element is expired it is removed from cache. Each time an item is fetched from the cache it is promoted to front provided it is not expired. When a new element is to be added and the cache has reached capacity, the least recently used item is removed from the cache.

## redisproxy

The redisproxy package manages the connection to the cache as well as the backing Redis service instance. The Get method will try to fetch a value from the LRU cache. If no value is found, it tries to get the value from the Redis server and if successful, adds it to the cache and returns the value back. The Redis proxy also implements a Redis client in GO that sends and receives RESP commands.

The Redis client is not a full featured redis client. I have implemented the basic redis send and receive so   values can be retrieved using this package. I used the stdlib bufio reader and writer to send the Redis commands and receive the response. The package can also process SET commands.  

# Assumptions

- The redis proxy server only currently supports GET of strings and can be expanded to handle the other data types as well as set. The underlying redis connection as part of the redisproxy package can handle set. It cannot handle complex responses like arrays yet.

- I have added some test data to the Redis instance to run tests. The assumption is that it will have data pre-propulated by another process.

- When a key is not found it will return a 404 with an empty value rather than an OK with a (nil) value.

## Pre-requisites
* Install Golang - https://golang.org/doc/install
* Set GOPATH and GOBIN

  ```
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$PATH
  ```
* Download and install [Docker for Mac](https://www.docker.com/products/docker#/mac)

## Config Environment variables Used

The server can be configured by setting environment variables
- SEGMENT_REDIS_SERVER
  Address of the backing Redis
- SEGMENT_REDIS_PORT
  Port of Redis server
- SEGMENT_CACHE_EXPIRY
  Cache expiry time in seconds

For docker these are configured in the docker compose file

## Quickstart guide

* Clone this repo into $GOPATH

* To run unit tests
    ```
    make test
    ```
* Run start script. If the API server is running, the script wont rebuild.
    ```
    ./start
    ```
* rebuild docker image with new code changes
    ```
    ./start rebuild
    ```
* stop script will stop docker containers
    ```
    scripts/docker/stop
    ```
* clean docker images so it can be rebuilt from sources
    ```
    scripts/docker/clean
    ```
* run redis-cli
  ```
  scripts/redis/redis-cli
  ```
* start redis without docker compose
    ```
    scripts/redis/start-redis
    ```

# References

- [Reddis](https://redis.io/commands/set)
- [Reddis Serialization Protocol] (https://redis.io/topics/protocol)
- [Golang](https://golang.org/pkg/)


# Contributors
* [Maria DeSouza](maria.g.desouza@gmail.com)
