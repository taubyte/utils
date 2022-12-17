package bundle

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Zip contents from source to target and returns the file.
// REF: https://forum.golangbridge.org/t/trying-to-zip-files-without-creating-folder-inside-archive/10260
func Zip(source, target string) (*os.File, error) {
	_, err := os.Stat(source)
	if err != nil {
		return nil, err
	}

	zipFile, err := os.Create(target)
	if err != nil {
		return nil, fmt.Errorf("Create `%s`: %w", target, err)
	}

	defer zipFile.Seek(0, io.SeekStart)
	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Walk basic source:`%s` failed with: %w", source, err)
		}

		if info.IsDir() {
			if source == path {
				return nil
			}
			path += "/"
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("Zipinfoheader `%s` failed with: %w", info, err)
		}

		header.Name = path[len(source)+1:]
		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("archive create header failed with: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Open Path:`%s` failed with: %w", path, err)
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("Copy failed with: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return zipFile, archive.Flush()
}
