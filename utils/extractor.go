package utils

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Extractor(src, tempdir string) error {

	supportFormat := []string{
		"rar",
		"zip",
		"tar",
		"gz",
		"br",
		"sz",
		"zz",
		"zst",
		"bz2",
		"7z",
		"xz",
		"tar",
		"tbz",
		"tar.xz",
		"tar.gz",
		"gzip",
	}

	var isSupported = false

	for _, format := range supportFormat {
		if strings.HasSuffix(src, format) {
			isSupported = true
			break // stop the loop and jump to the next instruction
		}

	}

	if !isSupported {
		// check if the file has no suffix at all
		if strings.IndexByte(src, '.') == -1 {
			// return fmt.Errorf("%s has no supported suffix", src)
			// just a binary, not an archive, or compressed archive
			return nil
			// do something else because it'a a regular file
		}
		return fmt.Errorf("%s is not supported", src)
	}

	reader, err := os.Open(src)
	// fullpath, _ := filepath.Abs(src)
	if err != nil {
		return err
	}

	format, input, err := archiver.Identify(src, reader)
	if err != nil {
		return err
	} else {
		if ex, ok := format.(archiver.Extractor); ok {
			// fmt.Println("Extracting ", src)
			ex.Extract(context.Background(), input, nil, func(ctx context.Context, f archiver.File) error {
				stat, _ := f.Stat()
				// fmt.Println(stat.Name())

				// create a new file with the same name as the extracted file
				// f.Open() returns an io.ReadCloser
				content, _ := f.Open()
				defer content.Close()

				newFilePath := filepath.Join(tempdir, stat.Name())
				newFile, err := os.Create(newFilePath)
				if err != nil {
					return err
				}
				defer newFile.Close()

				// copy the contents of the extracted file to the new file

				_, err = io.Copy(newFile, content)
				if err != nil {
					return err
				}
				return nil
			})

		}
	}
	return nil
}
