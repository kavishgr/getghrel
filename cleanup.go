package main

import (
	"debug/elf"
	"debug/macho"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

	err := filepath.Walk(tempdir, func(binpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a regular file
		if !info.Mode().IsRegular() {
			// fmt.Println("Not regular: ", binpath)
			return nil
		}

		// Open the file
		f, err := os.Open(binpath)
		defer f.Close()
		if err != nil {
			return nil
		}

		// Verify if the open file is either ELF or Mach-O
		if err := verifyFile(f); err == nil {
			err = os.Chmod(binpath, 0755)
			if err != nil {
				return err
			}
		} else {
			// fmt.Println("Removing: ", binpath)
			os.Remove(binpath)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return nil
}
