package algorithmia

/*
	IMPORTANT! For testing DataFile data interface use "test/datafile_test.go"
	           instead of this.
*/

import (
	"testing"
)

func TestDataFileSetAttributes(t *testing.T) {
	f := NewDataFile(NewClient("", ""), "data://a/b.txt")

	err := f.SetAttributes(&FileAttributes{
		LastModified: "2016-01-06T00:52:34.000Z",
		Size:         1,
	})

	if err != nil {
		t.Fatal(err)
	}

	ft := f.LastModified().Format("2006-01-02T15:04:05.000Z")
	if ft != "2016-01-06T00:52:34.000Z" {
		t.Fatal("got", ft, "after Format")
	}
}

func TestReadingErrorMessage(t *testing.T) {
	msg1 := []byte(`{"error":{"message":"everything is lost"}}`)

	e1, e2 := errFromJsonData(msg1)
	if e2 != nil {
		t.Fatal(e2)
	}

	if e1 == nil {
		t.Fatal("non-nil *Err expected")
	}

	if e1.Error() != "everything is lost" {
		t.Fatal("wrong error message")
	}
}
