package test

import (
	"os"
	"testing"

	algorithmia "github.com/algebraic-brain/algorithmia-go"
)

func TestAcl(t *testing.T) {
	const myPath = "data://.my/privatePermissions"

	client := algorithmia.NewClient(os.Getenv("ALGORITHMIA_API_KEY"), "")
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
