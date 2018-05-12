package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/mariadesouza/redisproxyserver/redisproxy"
)

/*


 */

type redisProxyServer struct {
	redis *redisproxy.Redisproxy
}

func main() {

	redisServer := os.Getenv("SEGMENT_REDIS_SERVER")
	if redisServer == "" {
		redisServer= "redis"
	}

	redisPort := os.Getenv("SEGMENT_REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	var cacheExpiry int64
	cacheExpiry, err := strconv.ParseInt(os.Getenv("SEGMENT_CACHE_EXPIRY"), 10, 64)
	if err != nil || cacheExpiry == 0 {
		cacheExpiry = 10
	}

	var redisProxy redisProxyServer

	redisProxy.redis, err = redisproxy.New(redisServer, redisPort, 10, cacheExpiry*60)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer redisProxy.redis.Close()
	http.HandleFunc("/", indexHandlerHelloWorld)
	http.HandleFunc("/redisproxy", redisProxy.IndexHandlerRedisproxy)
	log.Println("Redis Proxy webserver listening for requests")
	http.ListenAndServe(":8082", nil)
}

func indexHandlerHelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world, I'm running a redis proxy server on %s ", runtime.GOOS)
}

func (s *redisProxyServer) Close() {
	s.redis.Close()
}

func (s *redisProxyServer) IndexHandlerRedisproxy(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {
		case "GET":
			keys, ok := r.URL.Query()["key"]
			if (!ok || len(keys) < 1 ){
				 handleError(w, http.StatusBadRequest, "Url Param 'key' is missing")
				 return
			}
			key := keys[0]
			value, err := s.redis.Get(key)
			if err != nil {
				handleError(w, http.StatusNotFound, key+"not found:"+err.Error())
				return
			}
			response := fmt.Sprintf(`{"key": "%s","value": "%s"}`, key, value)
			fmt.Fprintln(w, response)
		case "PUT":
			 handleError(w, http.StatusNotImplemented,"" ) //TODO
		default:
			 handleError(w, http.StatusNotImplemented,"")
		}
	return

}


func handleError(w http.ResponseWriter, status int, message string) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	statusText := http.StatusText(status)

	eResp := fmt.Sprintf(`{"Error": "%s","Status": "%s"}`, message, statusText)

	fmt.Fprintln(w, eResp)
}
