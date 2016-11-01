package algorithmia

import (
	"fmt"
	"net/url"
	"strings"
)

type OutputType int

const (
	Default OutputType = iota
	Raw
	Void
)

type Algorithm struct {
	client *Client

	path            string
	url             string
	queryParameters url.Values
	outputType      OutputType
}

type AlgoOptions struct {
	Timeout         int
	Stdout          bool
	Output          OutputType
	QueryParameters url.Values
}

func NewAlgorithm(client *Client, ref string) (*Algorithm, error) {
	path := strings.TrimSpace(ref)

	if strings.HasPrefix(path, "algo:/") {
		path = path[len("algo:/"):]
	}

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	return &Algorithm{
		client:          client,
		path:            path,
		url:             "/v1/algo/" + path,
		queryParameters: url.Values{},
	}, nil
}

func (algo *Algorithm) SetOptions(opt AlgoOptions) {
	if opt.Timeout == 0 {
		opt.Timeout = 300
	}
	algo.queryParameters.Add("timeout", fmt.Sprint(opt.Timeout))
	algo.queryParameters.Add("stdout", fmt.Sprint(opt.Stdout)) // TODO: false? False? 0?
	algo.outputType = opt.Output
	if opt.QueryParameters != nil {
		for k, v := range opt.QueryParameters {
			for _, vv := range v {
				algo.queryParameters.Add(k, vv)
			}
		}
	}
}

func (algo *Algorithm) postRawOutput(input1 interface{}) ([]byte, error) {
	algo.queryParameters.Add("output", "raw")
	resp, err := algo.client.postJsonHelper(algo.url, input1, algo.queryParameters)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return getRaw(resp)
}

func (algo *Algorithm) postVoidOutput(input1 interface{}) (*AsyncResponse, error) {
	algo.queryParameters.Add("output", "void")
	resp, err := algo.client.postJsonHelper(algo.url, input1, algo.queryParameters)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return getAsyncResp(resp)
}

//Pipe an input into this algorithm
func (algo *Algorithm) Pipe(input1 interface{}) (interface{}, error) {
	switch algo.outputType {
	case Raw:
		return algo.postRawOutput(input1)
	case Void:
		return algo.postVoidOutput(input1)
	default:
		resp, err := algo.client.postJsonHelper(algo.url, input1, algo.queryParameters)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		return getAlgoResp(resp)
	}
}

func (algo *Algorithm) Path() string {
	return algo.path
}

func (algo *Algorithm) Url() string {
	return algo.url
}
