package algorithmia

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func getRaw(r *http.Response) ([]byte, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.New(string(b))
	}
	return b, nil
}

func getJson(r *http.Response, out interface{}) error {
	b, err := getRaw(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, out)
	return err
}

func getAsyncResp(r *http.Response) (*AsyncResponse, error) {
	var ar AsyncResponse
	err := getJson(r, &ar)
	if err != nil {
		return nil, err
	}

	if ar.Error != nil {
		return nil, ar.Error
	}

	return &ar, nil
}

func getAlgoResp(r *http.Response) (*AlgoResponse, error) {
	b, err := getRaw(r)
	if err != nil {
		return nil, err
	}

	return CreateAlgoResponse(b)
}
