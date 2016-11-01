package algorithmia

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type FileAttributes struct {
	FileName     string `json:"filename"`
	LastModified string `json:"last_modified"`
	Size         int64  `json:"size"`
}

type DataFile struct {
	DataObjectType

	path         string
	url          string
	lastModified time.Time
	size         int64

	client *Client
}

func NewDataFile(client *Client, dataUrl string) *DataFile {
	p := strings.TrimSpace(dataUrl)
	if strings.HasPrefix(p, "data://") {
		p = p[len("data://"):]
	} else if strings.HasPrefix(p, "/") {
		p = p[1:]
	}
	return &DataFile{
		DataObjectType: File,
		client:         client,
		path:           p,
		url:            getUrl(p),
	}
}

func (f *DataFile) SetAttributes(attr *FileAttributes) error {
	//%Y-%m-%dT%H:%M:%S.000Z
	t, err := time.Parse("2006-01-02T15:04:05.000Z", attr.LastModified)
	if err != nil {
		return err
	}
	f.lastModified = t
	f.size = attr.Size
	return nil
}

//Get file from the data api
func (f *DataFile) File() (*os.File, error) {
	if exists, err := f.Exists(); err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New(fmt.Sprint("file does not exist -", f.path))
	}

	resp, err := f.client.getHelper(f.url, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rf, err := ioutil.TempFile(os.TempDir(), "algorithmia")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(rf, resp.Body)
	if err != nil {
		rf.Close()
		return nil, err
	}
	name := rf.Name()
	if rf.Close(); err != nil {
		return nil, err
	}

	return os.Open(name)
}

func (f *DataFile) Exists() (bool, error) {
	resp, err := f.client.headHelper(f.url)
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, err
}

func (f *DataFile) Name() (string, error) {
	_, name, err := getParentAndBase(f.path)
	return name, err
}

func (f *DataFile) Bytes() ([]byte, error) {
	if exists, err := f.Exists(); err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New(fmt.Sprint("file does not exist -", f.path))
	}

	resp, err := f.client.getHelper(f.url, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (f *DataFile) StringContents() (string, error) {
	if b, err := f.Bytes(); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func (f *DataFile) Json(x interface{}) error {
	if exists, err := f.Exists(); err != nil {
		return err
	} else if !exists {
		return errors.New(fmt.Sprint("file does not exist -", f.path))
	}

	resp, err := f.client.getHelper(f.url, url.Values{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return getJson(resp, x)
}

//Post to data api
func (f *DataFile) Put(data interface{}) error {
	switch dt := data.(type) {
	case string:
		return f.PutBytes([]byte(dt))
	case []byte:
		return f.PutBytes(dt)
	default:
		return f.PutJson(data)
	}
}

func (f *DataFile) PutBytes(data []byte) error {
	resp, err := f.client.putHelper(f.url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := getRaw(resp)
	if err != nil {
		return err
	}

	err = errorFromJsonData(b)
	if err != nil {
		return err
	}

	return nil
}

//Post json to data api
func (f *DataFile) PutJson(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return f.PutBytes(b)
}

//Post file to data api
func (f *DataFile) PutFile(fpath string) error {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	return f.PutBytes(b)
}

//Delete from data api
func (f *DataFile) Delete() error {
	resp, err := f.client.deleteHelper(f.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := errorFromResponse(resp); err != nil {
		return err
	}

	return nil
}

func (f *DataFile) Path() string {
	return f.path
}

func (f *DataFile) Url() string {
	return f.url
}

func (f *DataFile) LastModified() time.Time {
	return f.lastModified
}

func (f *DataFile) Size() int64 {
	return f.size
}
