package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/nlopes/slack"
	"google.golang.org/appengine"
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
	// - - Answers (list questions and answers you've heard / made / found in deep dive explanation)
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

	client := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))

	http.HandleFunc("/raffler", raffleChallenge())
	http.HandleFunc("/raffler/start", slashCommand(raffleStart(client)))
	http.HandleFunc("/raffler/optin", slashCommand(raffleOptin(client)))
	http.HandleFunc("/raffler/optout", slashCommand(raffleOptout(client)))
	http.HandleFunc("/raffler/optinall", slashCommand(raffleOptinAll(client)))
	http.HandleFunc("/raffler/optoutall", slashCommand(raffleOptoutAll(client)))
	http.HandleFunc("/raffler/whosin", slashCommand(raffleWhosIn(client)))
	http.HandleFunc("/raffler/draw", slashCommand(raffleDraw(client)))
	http.HandleFunc("/raffler/stop", slashCommand(raffleStop(client)))
	http.HandleFunc("/raffler/set", slashCommand(raffleSetUsers(client)))
	http.HandleFunc("/game/connect", gameConnect)
	http.HandleFunc("/game/disconnect", gameConnect)
	http.HandleFunc("/game/listen", gameListen)
	http.HandleFunc("/game/send", gameSend)
	if os.Getenv("IN_APP_ENGINE") != "" {
		fmt.Println("Running in app engine")
		appengine.Main()
	} else {
		fmt.Println("Running on port 9092")
		http.ListenAndServe(":9092", nil)
	}
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
		err := tmpl.ExecuteTemplate(w, tmplName, inject)
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
	loadTemplates("templates")
}

func loadTemplates(container string) error {
	// Read the local templates directory to find .go.html files
	files, err := AssetDir(container)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, name := range files {
		// Can't use filepath.Ext because that splits at the last dot, and we
		// want the first dot.
		if !strings.HasSuffix(name, ".go.html") {
			continue
		}
		templates[name], err = template.ParseFiles(filepath.Join(container, name), filepath.Join(container, "footer.go.html"), filepath.Join(container, "header.go.html"))
		if err != nil {
			fmt.Println("Error decoding html template", name, ":", err)
		}
	}
	return nil
}
