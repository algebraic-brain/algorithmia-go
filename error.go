package algorithmia

type Err struct {
	Message    string `json:"message"`
	Stacktrace string `json:"stacktrace"`
}

func (e *Err) Error() string {
	return e.Message
}
