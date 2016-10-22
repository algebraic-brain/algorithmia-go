package algorithmia

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Metadata struct {
	ContentType string  `json:"content_type"`
	Duration    float64 `json:"duration"`
	Stdout      bool    `json:"stdout"`
}

func (d *Metadata) String() string {
	return fmt.Sprintf("Metadata(content_type=%v,duration=%v,stdout=%v)", d.ContentType, d.Duration, d.Stdout)
}

type AlgoResponse struct {
	Result   interface{} `json:"result"`
	Metadata *Metadata   `json:"metadata"`
	Error    *Err        `json:"error"` //never set!
}

func (resp *AlgoResponse) String() string {
	return fmt.Sprintf("AlgoResponse(result=%q,metadata=%v)", resp.Result, resp.Metadata)
}

func CreateAlgoResponse(b []byte) (*AlgoResponse, error) {
	var resp AlgoResponse

	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	if resp.Metadata.ContentType == "binary" {
		data, err := base64.StdEncoding.DecodeString(resp.Result.(string))
		if err != nil {
			return nil, err
		}
		resp.Result = data
	}

	return &resp, nil
}
