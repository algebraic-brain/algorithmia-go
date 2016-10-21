package algorithmia

import (
	"fmt"
)

type AsyncResponse struct {
	AsyncProtocol string `json:"async_protocol"`
	RequestId     string `json:"request_id"`
	Error         *Err   `json:"error"` //never set!
}

func (resp *AsyncResponse) String() string {
	return fmt.Sprint("AsyncResponse(async_protocol=%v, request_id=%v)", resp.AsyncProtocol, resp.RequestId)
}
