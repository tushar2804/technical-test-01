package main

import (
  "net/http"
  "net/http/httptest"
  "testing"
  "encoding/json"
)

func TestHealthCheckHandler(t *testing.T){
  req, err := http.NewRequest("GET", "/healthcheck", nil)
  if err != nil {
      t.Fatal(err)
  }

  rr := httptest.NewRecorder()
  handler := http.HandlerFunc(HealthCheckHandler)
  handler.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusOK {
      t.Errorf("handler returned wrong status code: got %v want %v",
          status, http.StatusOK)
  }

  healthStruct := HealthCheckInfo{}
  jsonErr := json.Unmarshal(rr.Body.Bytes(), &healthStruct)

  if jsonErr != nil {
    t.Errorf("Response can't be unmarshalled: %v", rr.Body.String())
  }

  if healthStruct.MyApplication[0].Description == "" {
    t.Errorf("No description in healthcheck.")
  }

  if healthStruct.MyApplication[0].LastCommitSHA == "" {
    t.Errorf("No commit SHA in healthcheck.")
  }

  if healthStruct.MyApplication[0].Version == "" {
    t.Errorf("No version in healthcheck.")
  }
}
