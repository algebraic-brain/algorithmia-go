package algorithmia

import (
	"errors"
	"path"
	"regexp"
	"strings"
)

var prefixRe = regexp.MustCompile(`([^:]+://)(/)?(.+)`)

var InvalidPath = errors.New("Invalid path")

func getParentAndBase(p string) (string, string, error) {
	var parent, base string
	if parts := prefixRe.FindStringSubmatch(p); parts == nil {
		var stripped string
		if strings.HasSuffix(p, "/") {
			stripped = p[:len(p)-1]
		} else {
			stripped = p
		}
		parent, base = path.Split(stripped)
		if parent == "" {
			return "", "", InvalidPath
		}
		parent = strings.TrimRight(parent, "/")
	} else {
		prefix, slash, uri := parts[1], parts[2], parts[3]
		parent, base = path.Split(uri)
		parent = strings.TrimRight(parent, "/")
		parent = prefix + slash + parent
	}
	return parent, base, nil
}

func PathJoin(parent, base string) string {
	if strings.HasSuffix(parent, "/") {
		return parent + base
	}
	return parent + "/" + base
}
