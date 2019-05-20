package main

import (
  "os"
  "os/exec"
  "net/http"
  "log"
  "encoding/json"
  "fmt"
  "time"
  "strings"

  "github.com/gorilla/mux"
)
// AppInfo has infomation for current build
type AppInfo struct {
  Version string `json:"version"`
  Description string `json:"description"`
  LastCommitSHA string `json:"lastcommitsha"`
}
// HealthCheckInfo has top level entry of myapplicatoin
type HealthCheckInfo struct {
  MyApplication []AppInfo `json:"myapplication"`
}

func getCommandOutput(cmd string, args ...string) string {
  stdout, err := exec.Command(cmd, args...).Output()
  if err == nil {
    return strings.TrimRight(string(stdout), "\n")
  }
  return err.Error()
}

// HealthCheckHandler handles http request to /healthcheck
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
  var CIVersion, CISHA, CIDescription string
  if os.Getenv("CI") == "" {
    // this is local dev mode
    CIVersion = "localdev"
    CIDescription = getCommandOutput("git", "log", "-1", "--pretty=%B")
    CISHA = getCommandOutput("git", "rev-parse", "HEAD")
  }else{
    CIVersion = os.Getenv("CI_VERSION")
    CIDescription = os.Getenv("CI_DESCRIPTION")
    CISHA = os.Getenv("CI_SHA")
  }
  jsonResp, _ := json.MarshalIndent(HealthCheckInfo{
    []AppInfo{
      {
        CIVersion,
    		CIDescription,
        CISHA,
      },
    },
	}, "", "    ")

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonResp))
}

func main() {
  port := os.Getenv("APP_PORT")
  if port == "" {
		port = "10000"
	}
  r := mux.NewRouter()
  r.HandleFunc("/healthcheck", HealthCheckHandler)
  server := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
  log.Fatal(server.ListenAndServe())
}
