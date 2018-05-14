# Redis Proxy webserver

A simple redis proxy server that uses an in memory LRU cache to speed up data access from redis. The cache items also have an expiration to keep data up to date.

# Design


## Server

The redisproxy server handles the GET request for a key. There are two packages implemented that are used by the server to implement the caching and redis connection management.


## lrucache

The LRU cache is implemented using a doubly linked list and a hashmap. It is guarded by a reader/writer mutual exclusion lock that will help prevent contention when high read rates occur concurrently. Each item has an expiration. If an element is expired it is removed from cache. Each time an item is fetched from the cache it is promoted to front provided it is not expired. When a new element is to be added and the cache has reached capacity, the least recently used item is removed from the cache.


## redisproxy

The redisproxy manages the connection to the redis server as well as to the cache. The GET method will try to fetch a value from the LRU cache. If no value is retrieved, it tries to get the value from the redis server and if successful adds it to the cache and returns it. The redis proxy also implements a client in GO.

The redis client is not a full featured redis client. I have implemented the basics so a GET and SET calls can be done using this package. I used the stdlib bufio reader and writer to send the Redis commands and receive the response. So far I have perfected the simple string and bulk string. I plan to do handle array responses as well if time permits.


# Assumptions

- The redis proxy server currently supports GET of strings and can be expanded to handle the other data types as well.

- The redis server will have  data pre-propulated by another process.


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

* Run start script
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
