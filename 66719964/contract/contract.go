package contract

type PingData struct {
	Name string
	Last time.Time
	Id   string
}

type PingRequestData struct {
	Name string
}

type BadInputResponse struct {
	status int
	Msg    string
	Schema interface{}
}

type PingUpdateResponse struct {
	Previous *PingData
	Current  PingData
}
