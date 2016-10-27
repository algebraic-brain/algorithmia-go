Algorithmia Common Library (Golang)
===================================

Golang client library for accessing the Algorithmia API
For API documentation, see the [Godoc](https://godoc.org/github.com/algebraic-brain/algorithmia-go)

## Install

```bash
go get github.com/algebraic-brain/algorithmia-go
```


## Authentication

First, create an Algorithmia client and authenticate with your API key:

```Go
import (
	algorithmia "github.com/algebraic-brain/algorithmia-go"
)

var apiKey = "{{Your API key here}}"
var client = algorithmia.NewClient(apiKey, "")
```

Now you're ready to call algorithms.

## Calling algorithms

The following examples of calling algorithms are organized by type of input/output which vary between algorithms.

Note: a single algorithm may have different input and output types, or accept multiple types of input,
so consult the algorithm's description for usage examples specific to that algorithm.

### Text input/output

Call an algorithm with text input by simply passing a string into its `pipe` method.
If the algorithm output is text, then the `result` field of the response will be a string.

```Go
	algo, _ := client.Algo("demo/Hello/0.1.1")
	response, _ := algo.Pipe("Author")
	fmt.Println(response.(*algorithmia.AlgoResponse).Result.(string))   //Hello Author
	fmt.Println(response.(*algorithmia.AlgoResponse).Metadata)          //Metadata(content_type=text,duration=0.0002127,stdout=false)
	fmt.Println(response.(*algorithmia.AlgoResponse).Metadata.Duration) //0.0002127
```

### JSON input/output

TODO
