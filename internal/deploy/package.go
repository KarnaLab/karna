package deploy

import (
	"os"

	"github.com/mholt/archiver"
)

func zipArchive(source, target string) (err error) {

	if _, err := os.Stat(target); err == nil {
		os.Remove(target)
	}

	if err = archiver.Archive([]string{source}, target); err != nil {
		return err
	}
	return
}
