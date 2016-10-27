package algorithmia

import (
	"errors"
	"fmt"
)

type aclInner struct {
	pseudonym string
	aclString string
}

func (a *aclInner) String() string {
	return fmt.Sprintf("AclType(pseudonym=%q,acl_string=%q)", a.pseudonym, a.aclString)
}

var aclTable = []*aclInner{
	&aclInner{"public", "user://*"},
	&aclInner{"my_algos", "algo://.my/*"},
	&aclInner{"private", ""},
}

type AclType int

func (t AclType) Pseudonym() string {
	return aclTable[t].pseudonym
}

func (t AclType) AclString() string {
	return aclTable[t].aclString
}

func (t AclType) String() string {
	return aclTable[t].String()
}

const (
	AclTypePublic AclType = iota
	AclTypeMyAlgos
	AclTypePrivate
	AclTypeDefault = AclTypeMyAlgos
)

var aclMap = map[string]AclType{
	AclTypePublic.AclString():  AclTypePublic,
	AclTypeMyAlgos.AclString(): AclTypeMyAlgos,
	AclTypePrivate.AclString(): AclTypePrivate,
}

func aclTypeFromResponse(aclList []string) (AclType, error) {
	if aclList == nil || len(aclList) == 0 {
		return AclTypePrivate, nil
	}

	if t, ok := aclMap[aclList[0]]; ok {
		return t, nil
	}

	return AclType(-1), errors.New(fmt.Sprint("Invalid acl string ", aclList[0]))
}

type Acl struct {
	readAcl AclType
}

func (a *Acl) ReadAcl() AclType {
	return a.readAcl
}

func (a *Acl) apiParam() *aclResponse {
	s := a.readAcl.AclString()
	if s == "" {
		return &aclResponse{Read: []string{}}
	}
	return &aclResponse{Read: []string{s}}
}

type aclResponse struct {
	Read []string `json:"read" mapstructure:"read"`
}

func (a *aclResponse) String() string {
	return fmt.Sprintf("AclResponse(read=%v)", a.Read)
}

var NoAclProvided = errors.New("Response does not contain read ACL")

func aclFromResponse(resp *aclResponse) (*Acl, error) {
	if resp.Read != nil {
		t, err := aclTypeFromResponse(resp.Read)
		if err != nil {
			return nil, err
		}
		return &Acl{readAcl: t}, nil
	}
	return nil, NoAclProvided
}

var ReadAclPublic = &Acl{AclTypePublic}
var ReadAclPrivate = &Acl{AclTypePrivate}
var ReadAclMyAlgos = &Acl{AclTypeMyAlgos}
