package contract

import(
  "time"
)

type PingData struct {
	Name string
	Last time.Time
	Id   string
}

type PingRequestData struct {
	Name string
}

type PingUpdateResponse struct {
	Previous *PingData
	Current  PingData
}

type BadInputResponse struct {
	Status int `json:omit`
	Msg    string
	Schema interface{}
}



