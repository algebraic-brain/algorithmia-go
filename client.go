//Algorithmia API Client (Go)
package algorithmia

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Client struct {
	apiKey     string
	apiAddress string
}

func NewClient(apiKey, apiAddress string) *Client {
	c := &Client{
		apiKey:     apiKey,
		apiAddress: apiAddress,
	}
	if apiAddress == "" {
		c.apiAddress = "https://api.algorithmia.com"
	}
	return c
}

func (c *Client) Algo(ref string) (*Algorithm, error) {
	return NewAlgorithm(c, ref)
}

func (c *Client) File(dataUrl string) *DataFile {
	return NewDataFile(c, dataUrl)
}

func (c *Client) Dir(dataUrl string) *DataDirectory {
	return NewDataDirectory(c, dataUrl)
}

func (c *Client) postJsonHelper(url string, input interface{}, params url.Values) (*http.Response, error) {
	headers := http.Header{}
	if c.apiKey != "" {
		headers.Add("Authorization", c.apiKey)
	}

	var (
		inputJson []byte
		err       error
	)
	if input == nil {
		headers.Add("Content-Type", "application/json")
		inputJson, err = json.Marshal(input)
		if err != nil {
			return nil, err
		}
	}

	switch inp := input.(type) {
	case string:
		headers.Add("Content-Type", "text/plain")
		inputJson = []byte(inp)
	case []byte:
		headers.Add("Content-Type", "application/octet-stream")
		inputJson = inp
	default:
		headers.Add("Content-Type", "application/json")
		inputJson, err = json.Marshal(input)
		if err != nil {
			return nil, err
		}
	}

	return request{Url: c.apiAddress + url, Data: inputJson, Headers: headers, Params: params}.post()
}

func (c *Client) getHelper(url string, params url.Values) (*http.Response, error) {
	headers := http.Header{}
	if c.apiKey != "" {
		headers.Add("Authorization", c.apiKey)
	}

	return request{Url: c.apiAddress + url, Headers: headers, Params: params}.get()
}

func (c *Client) headHelper(url string) (*http.Response, error) {
	headers := http.Header{}
	if c.apiKey != "" {
		headers.Add("Authorization", c.apiKey)
	}

	return request{Url: c.apiAddress + url, Headers: headers}.head()
}

func (c *Client) putHelper(url string, data []byte) (*http.Response, error) {
	headers := http.Header{}
	if c.apiKey != "" {
		headers.Add("Authorization", c.apiKey)
	}

	return request{Url: c.apiAddress + url, Headers: headers, Data: data}.put()
}

func (c *Client) deleteHelper(url string) (*http.Response, error) {
	headers := http.Header{}
	if c.apiKey != "" {
		headers.Add("Authorization", c.apiKey)
	}

	return request{Url: c.apiAddress + url, Headers: headers}.delete()
}

func (c *Client) patchHelper(url string, params map[string]interface{}) (*http.Response, error) {
	headers := http.Header{}
	headers.Add("content-type", "application/json")
	if c.apiKey != "" {
		headers.Add("Authorization", c.apiKey)
	}

	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return request{Url: c.apiAddress + url, Headers: headers, Data: b}.patch()
}
