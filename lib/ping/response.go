package ping

const StatusSuccess = 0
const StatusError = 1

// Response is the response return from ping probe
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Request *Request    `json:"request,omitempty"`
	Result  *Statistics `json:"result,omitempty"`
}
