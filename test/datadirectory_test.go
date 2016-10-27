package test

import (
	"os"
	"testing"

	algorithmia "github.com/algebraic-brain/algorithmia-go"
)

var client = algorithmia.NewClient(os.Getenv("ALGORITHMIA_API_KEY"), "")

func TestAcl(t *testing.T) {
	const myPath = "data://.my/privatePermissions"

	dd := client.Dir(myPath)
	defer dd.ForceDelete()

	if exists, err := dd.Exists(); err != nil {
		t.Fatal(err)

	} else if exists {
		if err := dd.ForceDelete(); err != nil {
			t.Fatal(err)
		}
	}

	if err := dd.Create(algorithmia.ReadAclPrivate); err != nil {
		t.Fatal(err)
	}

	perms, err := client.Dir(myPath).Permissions()
	if err != nil {
		t.Fatal(err)
	}

	if perms.ReadAcl() != algorithmia.AclTypePrivate {
		t.Fatal("private permissions expected")
	}

	err = dd.UpdatePermissions(algorithmia.ReadAclPublic)
	if err != nil {
		t.Fatal(err)
	}

	perms, err = client.Dir(myPath).Permissions()
	if err != nil {
		t.Fatal(err)
	}

	if perms.ReadAcl() != algorithmia.AclTypePublic {
		t.Fatal("public permissions expected, got", perms.ReadAcl())
	}
}

func TestDirName(t *testing.T) {
	dd := client.Dir("data://.my/this/is/a/long/path")
	n, err := dd.Name()
	if err != nil {
		t.Fatal(err)
	}
	if n != "path" {
		t.Fatal("'path' expected for directory name, got", n)

	}
}

func TestDirDoesNotExist(t *testing.T) {
	dd := client.Dir("data://.my/this_should_never_be_created")
	if exists, err := dd.Exists(); err != nil {
	} else if exists {
		t.Fatal("directory does not exist")
	}
}

func TestEmptyDirectoryCreationAndDeletion(t *testing.T) {
	/*
	   def test_empty_directory_creation_and_deletion(self):
	       dd = DataDirectory(self.client, "data://.my/empty_test_directory")

	       if (dd.exists()):
	           dd.delete(False)

	       self.assertFalse(dd.exists())

	       dd.create()
	       self.assertTrue(dd.exists())

	       # get rid of it
	       dd.delete(False)
	       self.assertFalse(dd.exists())
	*/

	dd := client.Dir("data://.my/empty_test_directory")

	if exists, err := dd.Exists(); err != nil {
		t.Fatal(err)
	} else if exists {
		if err := dd.Delete(); err != nil {
			t.Fatal(err)
		}
	}

	if exists, err := dd.Exists(); err != nil {
		t.Fatal(err)
	} else if exists {
		t.Fatal("directory does not exist")
	}

	err := dd.Create(nil)
	if err != nil {
		t.Fatal(err)
	}

	if exists, err := dd.Exists(); err != nil {
		t.Fatal(err)
	} else if !exists {
		t.Fatal("directory must exist")
	}

	if err := dd.Delete(); err != nil {
		t.Fatal(err)
	}

	if exists, err := dd.Exists(); err != nil {
		t.Fatal(err)
	} else if exists {
		t.Fatal("directory does not exist")
	}
}
