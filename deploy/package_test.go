package deploy

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestZipArchive(t *testing.T) {
	err := ioutil.WriteFile("./test", []byte{}, 0644)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = zipArchive("./test", "./lambda.zip")
	if err == nil {
		t.Log("Test PASSED because it build the archive with its content")
	} else {
		t.Errorf(err.Error())
	}

	os.Remove("./test")
	os.Remove("./lambda.zip")
}
