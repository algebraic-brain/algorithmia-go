package algorithmia

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func getUrl(p string) string {
	return "/v1/data/" + p
}

type datadirClient interface {
	getHelper(url string, params url.Values) (*http.Response, error)
	postJsonHelper(url string, input interface{}, params url.Values) (*http.Response, error)
	deleteHelper(url string) (*http.Response, error)
	patchHelper(url string, params map[string]interface{}) (*http.Response, error)
}

type DataDirectory struct {
	client datadirClient

	Path string
	Url  string
}

func NewDataDirectory(client datadirClient, dataUrl string) *DataDirectory {
	p := strings.TrimSpace(dataUrl)
	if strings.HasPrefix(p, "data://") {
		p = p[len("data://"):]
	} else if strings.HasPrefix(p, "/") {
		p = p[1:]
	}
	return &DataDirectory{
		client: client,
		Path:   p,
		Url:    getUrl(p),
	}
}

func (f *DataDirectory) Exists() (bool, error) {
	resp, err := f.client.getHelper(f.Url, url.Values{})
	return resp.StatusCode == http.StatusOK, err
}

func (f *DataDirectory) Name() (string, error) {
	_, name, err := getParentAndBase(f.Path)
	return name, err
}

func (f *DataDirectory) Create(acl *Acl) error {
	parent, name, err := getParentAndBase(f.Path)
	if err != nil {
		return err
	}
	jso := map[string]interface{}{
		"name": name,
	}
	if acl != nil {
		jso["acl"] = acl.ApiParam()
	}

	resp, err := f.client.postJsonHelper(getUrl(parent), jso, nil)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := getRaw(resp)
		if err != nil {
			return err
		}
		return errors.New("Directory creation failed: " + string(b))
	}
	return nil
}

func (f *DataDirectory) doDelete(force bool) error {
	url := f.Url
	if force {
		url += "?force=true"
	}

	resp, err := f.client.deleteHelper(url)
	if err != nil {
		return err
	}

	b, err := getRaw(resp)
	if err != nil {
		return err
	}

	err = ErrorFromJsonData(b)
	if err != nil {
		return err
	}

	return nil
}

func (f *DataDirectory) Delete() error {
	return f.doDelete(false)
}

func (f *DataDirectory) ForceDelete() error {
	return f.doDelete(true)
}

func (f *DataDirectory) File(name string) *DataFile {
	return NewDataFile(f.client.(*Client), PathJoin(f.Path, name))
}

func (f *DataDirectory) Dir(name string) *DataDirectory {
	return NewDataDirectory(f.client.(*Client), PathJoin(f.Path, name))
}
