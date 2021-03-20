package main

import(
  "github.com/farrellit/stackoverflow/66719964/contract"
  "net/http"
  "log"
  "bytes"
  "encoding/json"
  "sync"
)

const (
  surl = "http://localhost:9000"
)

var(
  client = &http.Client{}
)

func Create(name string) contract.PingData {
  var buf bytes.Buffer
  if err := json.NewEncoder(&buf).Encode(map[string]string{"Name": name}); err != nil {
    panic(err)
  }
  resp, err := http.Post(surl, "application/json", &buf)
  log.Println("Create: ", resp, err)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()
  var pd contract.PingData
  if err := json.NewDecoder(resp.Body).Decode(&pd); err != nil {
    panic(err)
  }
  return pd
}

func Get(id string) contract.PingData {
  resp, err := http.Get(surl + "/" + id)
  log.Println("Get: ", resp, err)
  if err != nil {
    panic(err)
  }
  var pd contract.PingData
  defer resp.Body.Close()
  if err := json.NewDecoder(resp.Body).Decode(&pd); err != nil {
    panic(err)
  }
  return pd
}

func Put(id, name string) contract.PingUpdateResponse {
  var buf bytes.Buffer
  if err := json.NewEncoder(&buf).Encode(map[string]string{"Name": name}); err != nil {
    panic(err)
  }
  req, err := http.NewRequest(http.MethodPut, surl + "/" + id, &buf)
  if err != nil {
    panic(err)
  }
  resp, err := client.Do(req)
  log.Println("Put: ", resp, err)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()
  var pd contract.PingUpdateResponse
  if err := json.NewDecoder(resp.Body).Decode(&pd); err != nil {
    panic(err)
  }
  return pd

}

func main() {
  var wg sync.WaitGroup
  for i := 0; i < 3; i++ {
    wg.Add(1)
    go dostuff(&wg)
  }
  wg.Wait()
}

func dostuff(wg *sync.WaitGroup) {
  defer wg.Done()
  pd := Create("df")
  log.Println(pd)

  for i := 0; i < 1000; i++ {

  pd = Get(pd.Id)
  log.Println(pd)

  pr := Put(pd.Id, "daniel")
  pd = pr.Current
  log.Println(pr,pd)

  pr = Put(pd.Id, "farrell")
  pd = pr.Current
  log.Println(pr,pd)

  pd = Get(pd.Id)
  log.Println(pd)
  }
}
