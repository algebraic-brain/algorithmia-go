package algorithmia

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Url     string
	Headers http.Header
	Params  url.Values
	Data    []byte
}

func (req Request) url() (string, error) {
	u, err := url.Parse(req.Url)

	if err != nil {
		return "", err
	}

	if req.Params != nil {
		u.RawQuery = req.Params.Encode()
	}

	return u.String(), nil
}

func (req Request) mkReq(method string, canData bool) (*http.Request, error) {
	u, err := req.url()
	if err != nil {
		return nil, err
	}

	var rdr io.Reader
	if canData && req.Data != nil {
		rdr = bytes.NewBuffer(req.Data)
	}

	r, err := http.NewRequest(method, u, rdr)
	if err != nil {
		return nil, err
	}

	if req.Headers != nil {
		for k, v := range req.Headers {
			if v == nil {
				continue
			}
			for _, vv := range v {
				r.Header.Add(k, vv)
			}
		}
	}

	return r, nil
}

func (req Request) Get() (*http.Response, error) {
	r, err := req.mkReq("GET", false)
	if err != nil {
		return nil, err
	}
	return (&http.Client{}).Do(r)
}

func (req Request) Post() (*http.Response, error) {
	r, err := req.mkReq("POST", true)
	if err != nil {
		return nil, err
	}
	return (&http.Client{}).Do(r)
}

func (req Request) Head() (*http.Response, error) {
	r, err := req.mkReq("HEAD", false)
	if err != nil {
		return nil, err
	}
	return (&http.Client{}).Do(r)
}

func (req Request) Put() (*http.Response, error) {
	r, err := req.mkReq("PUT", true)
	if err != nil {
		return nil, err
	}
	return (&http.Client{}).Do(r)
}

func (req Request) Delete() (*http.Response, error) {
	r, err := req.mkReq("DELETE", false)
	if err != nil {
		return nil, err
	}
	return (&http.Client{}).Do(r)
}

func (req Request) Patch() (*http.Response, error) {
	r, err := req.mkReq("PATCH", true)
	if err != nil {
		return nil, err
	}
	return (&http.Client{}).Do(r)
}
