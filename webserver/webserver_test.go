package main

import(
  "fmt"
  "testing"
  "net/http"
  "io/ioutil"
  "net/http/httptest"

  _"github.com/mariadesouza/redisproxyserver/redisproxy"
)


func TestErrorResponsesHandlerRedisproxy(t *testing.T) {

  var redisProxy redisProxyServer
  handler :=  redisProxy.HandlerRedisproxyRequest

  // Create a request to pass to the handler.
  badReq := httptest.NewRequest("GET", "/redisproxy", nil)
  w := httptest.NewRecorder()
  handler(w,badReq)
  resp := w.Result()
  if resp.StatusCode != http.StatusBadRequest{
    t.Errorf("Expected: ", http.StatusBadRequest, "Got:", resp.StatusCode)
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

func (f *mockHandler) HandlerRedisproxyRequest(w http.ResponseWriter, r *http.Request)  {
  keys, _ := r.URL.Query()["key"]
  var value string
  if len(keys) > 0{
    value = "pong"
  }
  response := fmt.Sprintf(`{"key": "%s","value": "%s"}`, keys[0], value)
  fmt.Fprintln(w, response)
  return

}

func TestSuccessResponsesHandlerRedisproxy(t *testing.T) {

  redisProxy := &mockHandler{}
  handler :=  redisProxy.HandlerRedisproxyRequest

  // Create a request to pass to the handler.-  mock here
  req := httptest.NewRequest("GET", "/redisproxy?key=ping", nil)
  w := httptest.NewRecorder()
  handler(w,req)
  resp := w.Result()
  if resp.StatusCode != http.StatusOK {
    t.Errorf("Expected: ", http.StatusBadRequest, "Got:", resp.StatusCode)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
  }
}
