package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	algorithmia "github.com/algebraic-brain/algorithmia-go"
)

func TestAlgo(t *testing.T) {
	c := algorithmia.NewClient(os.Getenv("ALGORITHMIA_API_KEY"), "")
	algo, err := c.Algo("demo/Hello")
	if err != nil {
		t.Fatal(err)
	}
	r, err := algo.Pipe(nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, ok := r.(*algorithmia.AlgoResponse)
	if !ok {
		t.Fatal("did not receive an AlgoResponse")
	}

	b, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}
