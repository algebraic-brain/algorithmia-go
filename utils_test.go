package algorithmia

import (
	"testing"
)

func TestGetParentAndBase(t *testing.T) {
	assertPB := func(pref, parentHave, parentMust, baseHave, baseMust string) {
		if parentHave != parentMust {
			t.Fatal(pref+"got parent", `"`+parentHave+`"`, "intstead of", `"`+parentMust+`"`)
		}
		if baseHave != baseMust {
			t.Fatal(pref+"got parent", `"`+baseHave+`"`, "intstead of", `"`+baseMust+`"`)
		}
	}
	p, b, err := getParentAndBase("a/b/c")
	if err != nil {
		t.Fatal(err)
	}
	assertPB("1:", p, "a/b", b, "c")

	p, b, err = getParentAndBase("data://foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	assertPB("2:", p, "data://foo", b, "bar")

	p, b, err = getParentAndBase("data:///foo")
	if err != nil {
		t.Fatal(err)
	}
	assertPB("3:", p, "data:///", b, "foo")

	p, b, err = getParentAndBase("data://foo")
	if err != nil {
		t.Fatal(err)
	}
	assertPB("4:", p, "data://", b, "foo")

	if _, _, err := getParentAndBase("/"); err == nil {
		t.Fatal("error expected")
	}

	if _, _, err := getParentAndBase(""); err == nil {
		t.Fatal("error expected")
	}

	if _, _, err := getParentAndBase("a/"); err == nil {
		t.Fatal("error expected")
	}
	if p := PathJoin("/a/b/c/", "d"); p != "/a/b/c/d" {
		t.Fatal(`"/a/b/c/d" expected, got`, p)
	}
	if p := PathJoin("/a/b/c", "d"); p != "/a/b/c/d" {
		t.Fatal(`"/a/b/c/d" expected, got`, p)
	}
	if p := PathJoin("/a//b/c//", "/d"); p != "/a//b/c///d" {
		t.Fatal(`"/a//b/c///d" expected, got`, p)
	}
}
