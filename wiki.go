package main

import (
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
)

const (
	DEFAULT_CONTENT_DIRECTORY string = "content"
	DEFAULT_PORT              int    = 1337
	DEFAULT_HTML_TEMPLATE     string = "view.html"
)

var (
	CONTENT_DIRECTORY string = DEFAULT_CONTENT_DIRECTORY
	PORT              int    = DEFAULT_PORT
	HTML_TEMPLATE     string = DEFAULT_HTML_TEMPLATE
	templates         *template.Template
)

type Page struct {
	Title string
	Body  template.HTML
}

func loadPage(page string) (*Page, error) {
	filename := "git/" + page + ".md"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: page, Body: template.HTML(blackfriday.MarkdownCommon([]byte(body)))}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path[1:]
	if len(page) == 0 {
		http.Redirect(w, r, "/Index", 301)
		return
	}
	p, err := loadPage(page)
	if err != nil {
		http.Redirect(w, r, "/", 301)
		return
	}
	renderTemplate(w, "view", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	flag.IntVar(&PORT, "p", DEFAULT_PORT, "Server port")
	flag.StringVar(&CONTENT_DIRECTORY, "C", DEFAULT_CONTENT_DIRECTORY, "Content directory")
	flag.StringVar(&HTML_TEMPLATE, "t", DEFAULT_HTML_TEMPLATE, "Html template")
	flag.Parse()

	templates = template.Must(template.ParseFiles(HTML_TEMPLATE))
}

func main() {
	http.HandleFunc("/", viewHandler)
	fmt.Printf("Magic happens on port %v...\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
	if nil != err {
		panic(err)
	}
}
