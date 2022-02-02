package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err := filepath.WalkDir("../quickstart", func(path string, d fs.DirEntry, err error) error {
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
	fmt.Println(err)
}
