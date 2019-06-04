package main

import (
	"net/http"

	"google.golang.org/appengine"

	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
)

type Page struct {
	Name     string
	Link     string
	SubPages []Page
}

func main() {
	// Site design:
	// Home
	// - Blog
	// - Games
	// - Contact
	// - Resume
	// - Projects
	// - - Oak
	// - - Geva

	pages := []Page{
		{
			Name: "Blog",
			Link: "blog",
		}, {
			Name: "Games",
			Link: "games",
		}, {
			Name: "Contact",
			Link: "contact",
		}, {
			Name: "Resume",
			Link: "resume",
		}, {
			Name: "Projects",
			SubPages: []Page{
				{
					Name: "Oak",
					Link: "projects/oak",
				},
				{
					Name: "Geva",
					Link: "projects/geva",
				},
			},
		},
	}
	http.HandleFunc("/blog", WriteTemplate(nil, "construction"))
	http.HandleFunc("/games", WriteTemplate(nil, "construction"))
	http.HandleFunc("/contact", WriteTemplate(nil, "construction"))
	http.HandleFunc("/resume", WriteTemplate(nil, "construction"))
	http.HandleFunc("/projects/oak", WriteTemplate(nil, "construction"))
	http.HandleFunc("/projects/geva", WriteTemplate(nil, "construction"))
	http.HandleFunc("/index", WriteTemplate(struct{ Pages []Page }{pages}, "home"))
	http.HandleFunc("/", LocalRedirect("index"))
	appengine.Main()
}

func LocalRedirect(path string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		newURL := &url.URL{}
		newURL.Scheme = req.URL.Scheme
		newURL.Host = req.URL.Host
		newURL.Path = path
		http.Redirect(w, req, newURL.String(), http.StatusMovedPermanently)
	}
}

func WriteTemplate(inject interface{}, tmplName string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		tmpl := templates[tmplName+".go.html"]
		err := tmpl.Execute(w, inject)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// We'd like to keep template files as .go.html files so they have nice syntax highlighting,
// so we need to store those files in a constant relative path from this directory. (As opposed
// to the alternative of storing the code as string consts)

// templates is the complete set of preloaded templates
var templates = make(map[string]*template.Template)

func init() {
	_, file, _, _ := runtime.Caller(0)
	container := filepath.Join(filepath.Dir(file), "templates")
	loadTemplates(container)
}

func loadTemplates(container string) error {
	// Read the local templates directory to find .go.html files
	fs, err := ioutil.ReadDir(container)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		// Can't use filepath.Ext because that splits at the last dot, and we
		// want the first dot.
		if !strings.HasSuffix(name, ".go.html") {
			continue
		}
		templates[name], err = template.ParseFiles(filepath.Join(container, name))
		if err != nil {
			fmt.Println("Error decoding html template", name, ":", err)
		}
	}
	return nil
}
