package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mariadesouza/cache/redisproxy"
)

func TestErrorResponsesHandlerRedisproxy(t *testing.T) {

	var redisProxy redisProxyServer
	handler := redisProxy.HandlerRedisproxyRequest

	// Create a request to pass to the handler.
	badReq := httptest.NewRequest("GET", "/redisproxy", nil)
	w := httptest.NewRecorder()
	handler(w, badReq)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected: %d Got: %d", http.StatusBadRequest, resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}

// This will mock the RedisProxy struct from the redisproxy package
// This is so our unit test doesnt have to make an actual connection to a redis instance
type RedisProxyConnector interface {
	HandlerRedisproxyRequest(w http.ResponseWriter, r *http.Request)
}

type mockHandler struct {
}

var _ RedisProxyConnector = (*mockHandler)(nil)

func (f *mockHandler) HandlerRedisproxyRequest(w http.ResponseWriter, r *http.Request) {
	keys, _ := r.URL.Query()["key"]
	var value string
	if len(keys) > 0 {
		value = "pong"
	}
	response := fmt.Sprintf(`{"key": "%s","value": "%s"}`, keys[0], value)
	fmt.Fprintln(w, response)
	return

}

func TestSuccessResponsesHandlerRedisproxy(t *testing.T) {

	redisProxy := &mockHandler{}
	handler := redisProxy.HandlerRedisproxyRequest

	// Create a request to pass to the handler.-  mock here

	req := httptest.NewRequest("GET", "/redisproxy?key=ping", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected: %d Got: %d", http.StatusBadRequest, resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}
