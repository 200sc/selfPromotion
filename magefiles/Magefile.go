//go:build mage

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh"
)

func Build() error {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	err := os.Chdir(filepath.Join(filepath.Dir(wd), "blog"))
	if err != nil {
		return err
	}
	err = sh.Run("hugo", "-D")
	if err != nil {
		return err
	}
	// I don't know why, but the index files that hugo generates don't work with a boring go http file server.
	// This renames them so that they do.
	return filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "index.html") {
			fmt.Println(path)
			newPath := filepath.Join(filepath.Dir(path), "index")
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			newFile, err := os.Create(newPath)
			if err != nil {
				return err
			}
			defer newFile.Close()
			_, err = newFile.Write(data)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func Deploy() error {
	// todo
	return nil
}
