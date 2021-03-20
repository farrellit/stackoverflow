package server

import (
  "github.com/farrellit/stackoverflow/66719964/contract"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
  "regexp"
)

type PingServer struct {
	data map[string]*contract.PingData
	lock sync.RWMutex
}

func NewPingServer() *PingServer {
	return &PingServer{
		data: make(map[string]*contract.PingData),
	}
}

func (ps *PingServer) Create(name string) (pd contract.PingData) {
	id := uuid.New().String()
	ps.lock.Lock()
	defer ps.lock.Unlock()
  ps.data[id] = &contract.PingData{Last: time.Now(), Name: name, Id: id}
  pd = *ps.data[id]
	return
}

func (ps *PingServer) Update(id, name string) (prev contract.PingData, found bool, cur contract.PingData) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
  var pprev *contract.PingData
	pprev, found = ps.data[id]
	if found {
    prev = *pprev
		if name != "" {
			ps.data[id].Name = name
		}
		ps.data[id].Last = time.Now()
		cur = *ps.data[id]
	}
	return
}

func (ps *PingServer) Read(id string) (pd contract.PingData, found bool) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	ppd, found := ps.data[id]
  if found {
    pd = *ppd
  }
  return
}

func (ps *PingServer) Delete(id string) (pd contract.PingData, found bool) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	ppd, found := ps.data[id]
  if found {
    pd = *ppd
    delete(ps.data, id)
  }
  return
}

func writeJson(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		e := fmt.Errorf("Error rendering %+v: %s", data, err.Error())
		panic(e)
	}
}

func (ps PingServerHttpController) WriteBadHttpResponse (
  br *contract.BadInputResponse,
  w http.ResponseWriter,
) {
	if br.Status == 0 {
		br.Status = http.StatusBadRequest
	}
	writeJson(w, br, br.Status)
}

type PingServerHttpController struct {
	ps *PingServer
}

func (ps *PingServer) HttpController() *PingServerHttpController {
	return &PingServerHttpController{ps: ps}
}

func (psc *PingServerHttpController) DecodeRequestOrRespond(w http.ResponseWriter, r *http.Request) (msg contract.PingRequestData, success bool) {
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		psc.WriteBadHttpResponse( &contract.BadInputResponse{
			  Msg:    err.Error(),
			  Schema: contract.PingRequestData{},
    }, w)
	} else {
		success = true
	}
	return
}

func (ps *PingServerHttpController) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch {
	case r.URL.Path == "/":
		ps.PostRequest(w, r)
	default:
		ps.RequestById(w, r)
	}
}

func (ps *PingServerHttpController) getIdFromPath(path string) (id, suffix string) {
	m := regexp.MustCompile("^/([^/]*)(.*)$").FindStringSubmatch(path)
	if m != nil {
		id = m[1]
		suffix = m[2]
	}
	return
}

func (ps *PingServerHttpController) RequestById(w http.ResponseWriter, r *http.Request) {
	id, _ := ps.getIdFromPath(r.URL.Path)
	notFound := func() {
		ps.WriteBadHttpResponse( &contract.BadInputResponse{
			Status: http.StatusNotFound,
			Msg:    "ID " + id + " not found",
			Schema: nil,
    }, w)
	}
  dataIfFound := func(pd contract.PingData, ok bool) {
		if !ok {
			notFound()
		} else {
			writeJson(w, pd, http.StatusOK)
		}
  }
	if id == "" {
		notFound()
		return
	}
	switch r.Method {
	case http.MethodGet:
		dataIfFound(ps.ps.Read(id))
	case http.MethodDelete:
		dataIfFound(ps.ps.Delete(id))
	case http.MethodPut:
		if msg, ok := ps.DecodeRequestOrRespond(w, r); ok {
			prev, isprev, cur := ps.ps.Update(id, msg.Name)
			if isprev {
				writeJson(w, contract.PingUpdateResponse{Previous: &prev, Current: cur}, http.StatusOK)
			} else {
				notFound()
			}
		}
	default:
		ps.WriteBadHttpResponse( &contract.BadInputResponse{
			Msg:    "This URL accepts only the GET or PUT method",
			Schema: contract.PingRequestData{},
			Status: http.StatusMethodNotAllowed,
		}, w )
	}
}

func (ps *PingServerHttpController) PostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ps.WriteBadHttpResponse(&contract.BadInputResponse{
			Msg:    "This URL accepts only the " + http.MethodPost + " Method",
			Schema: contract.PingRequestData{},
			Status: http.StatusMethodNotAllowed,
		}, w)
		return
	} else if msg, ok := ps.DecodeRequestOrRespond(w, r); ok {
		pd := ps.ps.Create(msg.Name)
		w.Header().Set("Location", "/"+pd.Id)
		writeJson(w, pd, http.StatusCreated)
	}
}

func main() {
	pingSvr := http.NewServeMux()
	pingSvr.HandleFunc("/", NewPingServer().HttpController().Handler)
	http.ListenAndServe(":9000", pingSvr)
}
