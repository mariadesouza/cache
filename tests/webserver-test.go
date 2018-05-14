package main

import (
  "net/http"
  "log"
  "encoding/json"
  "io/ioutil"
)

// This is basic functional testing of the redis proxy webserver
// Test data is added from scripts/redis/add-test-data

const (
	baseURL = "http://localhost:8082/"
)

var testTable = []struct {
  url string
  expectedCode int
  expectedResult string
} {
  {"redisproxy?key=jane", http.StatusOK, "mybuddy"},
  {"redisproxy?key=jack", http.StatusOK, "benimble"},
  {"redisproxy?key=unknown", http.StatusNotFound, ""},
}

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {


  for _, test := range testTable {
		req, err := http.NewRequest("GET", baseURL+test.url, nil)
    if err != nil {
			log.Println(err)
			continue
		}
    client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}
		// Check the status code is what we expect.
		if resp.StatusCode != test.expectedCode {
      log.Printf("%s returned wrong status code: got %v want %v\n", test.url, resp.StatusCode, test.expectedCode)
			continue
    }
    if resp.StatusCode == http.StatusOK {
			data, err := ioutil.ReadAll(resp.Body)
      var response Response
      err = json.Unmarshal(data, &response)
      if err != nil {
				log.Println(err)
				continue
			}
      if response.Value!= test.expectedResult{
        log.Printf("%s returned wrong response: got %s want %s\n", test.url, response.Value, test.expectedResult)
        continue
      }
      log.Println("OK: ", baseURL+test.url)
    } else if resp.StatusCode == test.expectedCode {
      log.Printf("OK: %s returned expected code %d\n", baseURL+test.url, test.expectedCode)
    }

  }
}
