package deploy

import (
	"karna/core"
	"os"

	"github.com/mholt/archiver"
)

func zipArchive(source, target string) {

	if _, err := os.Stat(target); err == nil {
		os.Remove(target)
	}

	err := archiver.Archive([]string{source}, target)

	if err != nil {
		core.LogErrorMessage(err.Error())

	}
}
