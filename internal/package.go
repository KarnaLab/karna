package deploy

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	archiver "github.com/mholt/archiver/v3"
)

func zipArchive(source, target, outputPathWithoutArchive string, isExecutable bool) (err error) {

	if _, err := os.Stat(target); err == nil {
		os.Remove(target)
	}

	if isExecutable {
		if err := compressExeAndArgs(target, source, outputPathWithoutArchive); err != nil {
			return fmt.Errorf("failed to compress file: %v", err)
		}
	} else {
		if err = archiver.Archive([]string{source}, target); err != nil {
			return err
		}
	}

	return
}

func writeExe(writer *zip.Writer, pathInZip string, data []byte) error {
	if pathInZip != "bootstrap" {
		header := &zip.FileHeader{Name: "bootstrap", Method: zip.Deflate}
		header.SetMode(0755 | os.ModeSymlink)

		link, err := writer.CreateHeader(header)

		if err != nil {
			return err
		}
		if _, err := link.Write([]byte(pathInZip)); err != nil {
			return err
		}
	}

	exe, err := writer.CreateHeader(&zip.FileHeader{
		CreatorVersion: 3 << 8,     // indicates Unix
		ExternalAttrs:  0777 << 16, // -rwxrwxrwx file permissions
		Name:           pathInZip,
		Method:         zip.Deflate,
	})
	if err != nil {
		return err
	}

	_, err = exe.Write(data)
	return err
}

func compressExeAndArgs(outZipPath, exePath, outputPathWithoutArchive string) error {
	err := os.MkdirAll(outputPathWithoutArchive, os.ModePerm)

	if err != nil {
		return err
	}

	zipFile, err := os.Create(outZipPath)

	if err != nil {
		return err
	}

	defer func() {
		closeErr := zipFile.Close()

		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close zip file: %v\n", closeErr)
		}
	}()

	zipWriter := zip.NewWriter(zipFile)

	defer zipWriter.Close()

	data, err := ioutil.ReadFile(exePath)

	if err != nil {
		return err
	}

	err = writeExe(zipWriter, filepath.Base(exePath), data)

	if err != nil {
		return err
	}

	return err
}
