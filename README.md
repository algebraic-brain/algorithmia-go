Algorithmia Common Library (Golang)
===================================

Golang client library for accessing the Algorithmia API.

For API documentation, see the [Godoc](https://godoc.org/github.com/algebraic-brain/algorithmia-go).

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

Call an algorithm with text input by simply passing a string into its `Pipe` method.
If the algorithm output is text, then the `Result` field of the response will be a string.

```Go
algo, _ := client.Algo("demo/Hello/0.1.1")
resp, _ := algo.Pipe("Author")
response := resp.(*algorithmia.AlgoResponse)
fmt.Println(response.Result)            //Hello Author
fmt.Println(response.Metadata)          //Metadata(content_type='text',duration=0.0002127)
fmt.Println(response.Metadata.Duration) //0.0002127
```

### JSON input/output

Call an algorithm with JSON input by simply passing in a type that can be serialized to JSON.
For algorithms that return JSON, the `Result` field of the response will be the appropriate
deserialized type.

```Go
algo, _ := client.Algo("WebPredict/ListAnagrams/0.1.0")
resp, _ := algo.Pipe([]string{"transformer", "terraforms", "retransform"})
response := resp.(*algorithmia.AlgoResponse)
fmt.Println(response.Result) //[transformer retransform]
```

### Binary input/output

Call an algorithm with binary input by passing a byte array into the `Pipe` method.
Similarly, if the algorithm response is binary data, then the `Result` field of the response
will be a byte array.

```Go
input, _ := ioutil.ReadFile("/path/to/bender.png")
algo, _ := client.Algo("opencv/SmartThumbnail/0.1")
resp, _ := algo.Pipe(input)
response := resp.(*algorithmia.AlgoResponse)
ioutil.WriteFile("thumbnail.png", response.Result.([]byte), 0666)
fmt.Println(response.Result) //[binary byte sequence]
```

### Error handling

API errors and Algorithm exceptions will result in calls to `Pipe` returning an error:

```Go
algo, _ := client.Algo("util/whoopsWrongAlgo")
_, err := algo.Pipe("Hello, World!")
fmt.Println(err) //algorithm algo://util/whoopsWrongAlgo not found
```

### Request options

The client exposes options that can configure algorithm requests.
This includes support for changing the timeout or indicating that the API should include stdout in the response.

```Go
algo, _ = client.Algo("util/echo")
algo.SetOptions(algorithmia.AlgoOptions{Timeout: 60, Stdout: false})
```

## Working with data

### Create directories
Create directories by instantiating a `DataDirectory` object and calling `Create`:

```Go
client.Dir("data://.my/foo").Create(nil) //nil for default access control (private)
```

### Upload files to a directory

Upload files by calling `Put` on a `DataFile` object.

```Go
foo := client.Dir("data://.my/foo")
foo.File("sample.txt").Put("sample text contents")
foo.File("binary_file").Put([]byte{72, 101, 108, 108, 111})
```

Note: you can instantiate a `DataFile` by either `client.File(filepath)` or `client.Dir(path).File(filename)`

### Download contents of file

Download files by calling `StringContents`, `Bytes`, `Json`, or `File` on a `DataFile` object:

```Go
foo := client.Dir("data://.my/foo")
sampleText, _ := foo.File("sample.txt").StringContents() //string object
fmt.Println(sampleText)                                  //"sample text contents"

binaryContent, _ := foo.File("binary_file").Bytes()      //binary data
fmt.Println(string(binaryContent))                       //"Hello"

tempFile, _ := foo.File("binary_file").File()            //Open file descriptor for read
defer tempFile.Close()
binaryContent, _ = ioutil.ReadAll(tempFile)
fmt.Println(string(binaryContent))                       //"Hello"
```

### Delete files and directories

Delete files and directories by calling `Delete` on their respective `DataFile` or `DataDirectory` object.
DataDirectories have `ForceDelete` method that deletes the directory even it contains files or other directories.

```Go
foo := client.Dir("data://.my/foo")
foo.File("sample.txt").Delete()
foo.ForceDelete() // force deleting the directory and its contents
```

### List directory contents

Iterate over the contents of a directory using the channel returned by calling `List`, `Files`, or `Dirs`
on a `DataDirectory` object:

```Go
foo := client.Dir("data://.my/foo")

// List files in "foo"
for entry := range foo.Files() {
	if entry.Err == nil {
		file := entry.Object.(*algorithmia.DataFile)
		fmt.Println(file.Path(), "at URL:", file.Url(), "last modified:", file.LastModified())
	}
}

// List directories in "foo"
for entry := range foo.Dirs() {
	if entry.Err == nil {
		dir := entry.Object.(*algorithmia.DataDirectory)
		fmt.Println(dir.Path(), "at URL:", dir.Url())
	}
}

// List everything in "foo"
for entry := range foo.List() {
	if entry.Err == nil {
		fmt.Println(entry.Object.Path(), "at URL:", entry.Object.Url())
	}
}
```

### Manage directory permissions

Directory permissions may be set when creating a directory, or may be updated on already existing directories.

```Go
foo := client.Dir("data://.my/foo")

//ReadAclPublic is a wrapper for &Acl{AclTypePublic} to make things easier
foo.Create(algorithmia.ReadAclPublic)

acl, _ := foo.Permissions()                             //Acl object
fmt.Println(acl.ReadAcl() == algorithmia.AclTypePublic) //true

foo.UpdatePermissions(algorithmia.ReadAclPrivate)
acl, _ = foo.Permissions()                               //Acl object
fmt.Println(acl.ReadAcl() == algorithmia.AclTypePrivate) //true
```

## Running tests

To run all test files:
```bash
export ALGORITHMIA_API_KEY={{Your API key here}}
cd test
go test -v
```

To run particular test:
```bash
export ALGORITHMIA_API_KEY={{Your API key here}}
cd test -v
go test datadirlarge_test.go -v
```
