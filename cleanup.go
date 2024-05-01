package main

import (
	"debug/elf"
	"debug/macho"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"io/fs"
)

func cleanup(tempdir string) error {
	var verifyFile func(file *os.File) error

	switch runtime.GOOS {
	case "linux":
		verifyFile = func(file *os.File) error {
			_, err := elf.NewFile(file)
			return err
		}
	case "darwin":
		verifyFile = func(file *os.File) error {
			_, err := macho.NewFile(file)
			return err
		}
	}

	err := filepath.WalkDir(tempdir, func(binpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if binpath == tempdir {
      		return nil
    	}

		// If a directory - skip
		if d.IsDir() {
			// fmt.Println("Not regular: ", binpath) // comment
			return nil
		}

		// Open the file
		f, err := os.Open(binpath)
		defer f.Close()
		if err != nil {
			return err
		}

		// Verify if the open file is either an ELF or Mach-O
		if err := verifyFile(f); err == nil {
			err = os.Chmod(binpath, 0755)
			// fmt.Println("Binary: ", binpath) //comment
			if err != nil {
				return err
			}
		} else {
			// fmt.Println("Removing: ", binpath) // comment
			os.Remove(binpath)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return nil
}
