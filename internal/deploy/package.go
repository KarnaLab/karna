package deploy

import (
	"github.com/mholt/archiver"
)

func zipArchive(source, target string) {

	err := archiver.Archive([]string{source}, target)

	if err != nil {
		panic(err.Error())
	}
}
