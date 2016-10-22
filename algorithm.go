package algorithmia

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type OutputType int

const (
	Default OutputType = iota
	Raw
	Void
)

type algoClient interface {
	postJsonHelper(url string, input interface{}, params url.Values) (*http.Response, error)
}

type Algorithm struct {
	Client          algoClient
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

func NewAlgorithm(client algoClient, ref string) (*Algorithm, error) {
	path := strings.TrimSpace(ref)

	if strings.HasPrefix(path, "algo:/") {
		path = path[len("algo:/"):]
	}

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	return &Algorithm{
		Client:          client,
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
	resp, err := algo.Client.postJsonHelper(algo.Url, input1, algo.QueryParameters)
	if err != nil {
		return nil, err
	}

	return getRaw(resp)
}

func (algo *Algorithm) postVoidOutput(input1 interface{}) (*AsyncResponse, error) {
	algo.QueryParameters.Add("output", "void")
	resp, err := algo.Client.postJsonHelper(algo.Url, input1, algo.QueryParameters)
	if err != nil {
		return nil, err
	}

	return getAsyncResp(resp)
}

func (algo *Algorithm) Pipe(input1 interface{}) (interface{}, error) {
	switch algo.OutputType {
	case Raw:
		return algo.postRawOutput(input1)
	case Void:
		return algo.postVoidOutput(input1)
	default:
		resp, err := algo.Client.postJsonHelper(algo.Url, input1, algo.QueryParameters)
		if err != nil {
			return nil, err
		}

		return getAlgoResp(resp)
	}
}
