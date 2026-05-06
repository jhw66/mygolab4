package serializer

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Msg    string      `json:"msg,omitempty"`
	Error  string      `json:"error,omitempty"`
}
