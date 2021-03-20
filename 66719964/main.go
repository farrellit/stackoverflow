package main

import (
  "github.com/farrellit/stackoverflow/66719964/server"
	"net/http"
)


func main() {
	pingSvr := http.NewServeMux()
	pingSvr.HandleFunc("/", server.NewPingServer().HttpController().Handler)
	err := http.ListenAndServe(":9000", pingSvr)
  panic(err)
}
