package bundle

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Options declares overrides to be made to files while creating the tarball.
type Options struct {
	FileOptions
}

// FileOptions declares changes to be made to individual files when walking a directory
type FileOptions struct {
	AccessTime time.Time
	ChangeTime time.Time
	ModTime    time.Time
}

func Tarball(src string, options *Options, writers ...io.Writer) error {
	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("unable to find source with %s", err)
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		// return on any error
		if err != nil {
			return err
		}

		// return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// This way bytes read is always same as long as data is same
		if options != nil {
			if options.AccessTime.IsZero() == false {
				header.AccessTime = options.AccessTime
			}
			if options.ChangeTime.IsZero() == false {
				header.ChangeTime = options.ChangeTime
			}
			if options.ModTime.IsZero() == false {
				header.ModTime = options.ModTime
			}
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		_, err = io.Copy(tw, f)

		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()
		if err != nil {
			return err
		}

		return nil
	})
}
