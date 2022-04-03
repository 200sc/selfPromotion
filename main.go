package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"google.golang.org/appengine"
)

func main() {
	manualPath := http.Dir(filepath.Join(".", "blog", "public"))

	httpfs := http.FileServer(manualPath)
	http.Handle("/", httpfs)

	if os.Getenv("IN_APP_ENGINE") != "" {
		fmt.Println("Running in app engine")
		appengine.Main()
	} else {
		fmt.Println("Running on port 9092")
		http.ListenAndServe(":9092", nil)
	}
}
