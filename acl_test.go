package algorithmia

import (
	"testing"
)

func TestAcl(t *testing.T) {
	if s := AclTypePrivate.AclString(); s != "" {
		t.Fatal(int(AclTypePrivate), "empty string expected, got", s)
	}

	if s := AclTypeMyAlgos.AclString(); s != "algo://.my/*" {
		t.Fatal(int(AclTypeMyAlgos), "'algo://.my/*' expected, got", s)
	}

	if s := AclTypePublic.AclString(); s != "user://*" {
		t.Fatal(int(AclTypePublic), "'user://*' expected, got", s)
	}

	if s := AclTypeDefault.AclString(); s != AclTypeMyAlgos.AclString() {
		t.Fatal(int(AclTypeDefault), AclTypeMyAlgos.AclString()+" expected, got", s)
	}

	if typ, err := AclTypeFromResponse([]string{}); err != nil {
		t.Fatal(err)
	} else if typ != AclTypePrivate {
		t.Fatal("AclTypePrivate expected, got", typ)
	}

	if typ, err := AclTypeFromResponse([]string{}); err != nil {
		t.Fatal(err)
	} else if typ != AclTypePrivate {
		t.Fatal("AclTypePrivate expected, got", typ)
	}

	if typ, err := AclTypeFromResponse([]string{"algo://.my/*"}); err != nil {
		t.Fatal(err)
	} else if typ != AclTypeMyAlgos {
		t.Fatal("AclTypeMyAlgos expected, got", typ)
	}

	if typ, err := AclTypeFromResponse([]string{"user://*"}); err != nil {
		t.Fatal(err)
	} else if typ != AclTypePublic {
		t.Fatal("AclTypePublic expected, got", typ)
	}
}
