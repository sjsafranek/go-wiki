package main

import (
	// "errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/russross/blackfriday"
)

const (
	DEFAULT_CONTENT_DIRECTORY       string = "content"
	DEFAULT_HTML_TEMPLATE_FILE_FILE string = "view.html"
)

var (
	CONTENT_DIRECTORY  string = DEFAULT_CONTENT_DIRECTORY
	HTML_TEMPLATE_FILE string = DEFAULT_HTML_TEMPLATE_FILE_FILE
	HTML_TEMPLATE_NAME string = ""
	TEMPLATES          *template.Template
)

type Page struct {
	Title string
	Body  template.HTML
}

type WikiEngine struct{}

func (self *WikiEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	self.viewHandler(w, r)
}

func (self *WikiEngine) getFilename(page string) string {
	return fmt.Sprintf("%v/%v.md", CONTENT_DIRECTORY, page)
}

func (self *WikiEngine) loadPage(page string) (*Page, error) {
	filename := self.getFilename(page)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: page, Body: template.HTML(blackfriday.MarkdownCommon([]byte(body)))}, nil
}

func (self *WikiEngine) viewHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path[1:]
	if len(page) == 0 {
		page = "index"
	}
	p, err := self.loadPage(page)
	if err != nil {
		self.renderTemplate(w, HTML_TEMPLATE_NAME, &Page{})
		return
	}
	self.renderTemplate(w, HTML_TEMPLATE_NAME, p)
}

func (self *WikiEngine) renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := TEMPLATES.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
