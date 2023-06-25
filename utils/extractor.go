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
		"bz2",
		"7z",
		"xz",
		"tar.xz",
		"tar.gz",
		"gzip",
	}

	var isSupported = false
	for _, format := range supportFormat {
		if strings.HasSuffix(src, format) {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return fmt.Errorf("%s is not supported", src)
	} else {
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
	}
	return nil
}
