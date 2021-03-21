package main

import (
  "github.com/farrellit/stackoverflow/66719964/server"
	"net/http"
  "os"
  "log"
  "runtime/pprof"
)


func main() {
  f, err := os.Create("main.prof")
  if err != nil {
    log.Fatal(err)
  }
  pprof.StartCPUProfile(f)
  defer pprof.StopCPUProfile()
	pingSvr := http.NewServeMux()
	pingSvr.HandleFunc("/", server.NewPingServer().HttpController().Handler)
	if err := http.ListenAndServe(":9000", pingSvr); err != nil {
    panic(err)
  }
}
