package main

import(
  "github.com/farrellit/stackoverflow/66719964/contract"
  "net/http"
  "log"
  "bytes"
  "encoding/json"
  "sync"
  "time"
  "fmt"
)

const (
  surl = "http://127.0.0.1:9000"
)

type PingClient struct {
  client *http.Client
  id int
  backoff int
  name string
  pd *contract.PingData
}

func (pc *PingClient)Backoff(err error) {
  if pc.backoff == 0 {
    pc.backoff = 1000
    return
  }
  pc.backoff = pc.backoff * 2
  if pc.backoff > 20000 {
    pc.backoff = 20000
  }
   log.Printf("Client %d: Backing off %d ms: %s", pc.id, pc.backoff, err.Error())
  time.Sleep(time.Duration(pc.backoff)*time.Millisecond)
}

func Decode(resp *http.Response, target interface{}) {
  defer resp.Body.Close()
  if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
    panic(err)
  }
}

func (pc *PingClient)Encode(data interface{}) *bytes.Buffer {
  var buf = new(bytes.Buffer)
  if err := json.NewEncoder(buf).Encode(data); err != nil {
    panic(err)
  }
  return buf
}

func (pc *PingClient)Create() {
  buf := pc.Encode(map[string]string{"Name": pc.name})
  req , err := http.NewRequest(http.MethodPost, surl, buf)
  if err != nil {
    panic(err)
  }
  pc.DoRequest(req, pc.decodeToData)
}

func (pc *PingClient)DoRequest(req *http.Request, success func(*http.Response)){
  req.Header.Set("Content-Type", "application/json")
  for {
    resp, err := pc.client.Do(req)
    if err != nil {
      pc.Backoff(err)
      continue
    }
    pc.backoff = 0
    if resp.StatusCode == 404 {
      panic(fmt.Errorf("Request failed with 404 for %s", req.URL.Path))
    }
    success(resp)
    break
  }
}

func (pc *PingClient)decodeToData(resp *http.Response){
    if pc.pd == nil {
      pc.pd = new(contract.PingData)
    }
    Decode(resp, &pc.pd)
}

func (pc *PingClient)Get() {
  req , err := http.NewRequest(http.MethodGet, surl + "/" + pc.pd.Id, nil)
  if err != nil {
    panic(err)
  }
  pc.DoRequest(req, pc.decodeToData)
}

func (pc *PingClient)Delete() {
  req, err := http.NewRequest(http.MethodDelete, surl + "/" + pc.pd.Id, nil)
  if err != nil {
    panic(err)
  }
  pc.DoRequest(req, func(resp *http.Response){
    pc.decodeToData(resp)
    log.Printf("Removed: %+v", pc.pd)
    pc.pd = nil
  })
}

func (pc *PingClient)Put(name string) {
  buf := pc.Encode(map[string]string{"Name": name})
  req, err := http.NewRequest(http.MethodPut, surl + "/" + pc.pd.Id, buf)
  if err != nil {
    panic(err)
  }
  pc.DoRequest(req, func(resp *http.Response){
    var pd contract.PingUpdateResponse
    Decode(resp, &pd)
    *pc.pd = pd.Current
    pc.name = pc.pd.Name
  })
}

func main() {
  var wg sync.WaitGroup
  client := &http.Client{}
  for i := 0; i < 50; i++ {
    wg.Add(1)
    go dostuff(client, i, &wg)
  }
  wg.Wait()
}

func dostuff(client *http.Client, id int, wg *sync.WaitGroup) {
  defer wg.Done()
  pc := PingClient{
    client: client,
    name: "df",
    id: id,
  }
  pc.Create()
  for i := 0; i < 1000; i++ {
    pc.Get()
    pc.Put("daniel")
    pc.Put("farrell")
    pc.Get()
  }
  pc.Delete()
}
