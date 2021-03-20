package main

import(
  "net/http"
  "log"
  "bytes"
  "encoding/json"
)

func main() {
  url := "http://localhost:9000"
  var buf bytes.Buffer
  if err := json.NewEncoder(&buf).Encode(map[string]string{"Name": "daniel"}); err != nil {
    panic(err)
  }
  resp, err := http.Post(url, "application/json", &buf)
  log.Println(resp,err)
  defer resp.Body.Close()
  var pd = make(map[string]string)
  if err := json.NewDecoder(resp.Body).Decode(&pd); err != nil {
    panic(err)
  }
  id, ok := pd["Id"]
  if !ok {
    panic("Expected Id field not in response")
  }
  log.Printf("%s - %+v", pd, id)
}
