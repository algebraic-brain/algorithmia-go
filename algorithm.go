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

	Path            string
	Url             string
	QueryParameters url.Values
	OutputType      OutputType
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
		Path:            path,
		Url:             "/v1/algo/" + path,
		QueryParameters: url.Values{},
	}, nil
}

func (algo *Algorithm) SetOptions(opt AlgoOptions) {
	algo.QueryParameters = opt.QueryParameters
	algo.QueryParameters.Add("timeout", fmt.Sprint(opt.Timeout))
	algo.QueryParameters.Add("stdout", fmt.Sprint(opt.Stdout)) // TODO: false? False? 0?
	algo.OutputType = opt.Output
}

func (algo *Algorithm) postRawOutput(input1 interface{}) ([]byte, error) {
	algo.QueryParameters.Add("output", "raw")
	resp, err := algo.client.postJsonHelper(algo.Url, input1, algo.QueryParameters)
	if err != nil {
		return nil, err
	}

	return getRaw(resp)
}

func (algo *Algorithm) postVoidOutput(input1 interface{}) (*AsyncResponse, error) {
	algo.QueryParameters.Add("output", "void")
	resp, err := algo.client.postJsonHelper(algo.Url, input1, algo.QueryParameters)
	if err != nil {
		return nil, err
	}

	return getAsyncResp(resp)
}

//Pipe an input into this algorithm
func (algo *Algorithm) Pipe(input1 interface{}) (interface{}, error) {
	switch algo.OutputType {
	case Raw:
		return algo.postRawOutput(input1)
	case Void:
		return algo.postVoidOutput(input1)
	default:
		resp, err := algo.client.postJsonHelper(algo.Url, input1, algo.QueryParameters)
		if err != nil {
			return nil, err
		}

		return getAlgoResp(resp)
	}
}
