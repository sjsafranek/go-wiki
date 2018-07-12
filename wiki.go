package main

import (
	// "errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/russross/blackfriday"
)

const (
	DEFAULT_WIKI_DIRECTORY          string = "wiki"
	DEFAULT_CONTENT_DIRECTORY       string = "content"
	DEFAULT_HTML_TEMPLATE_FILE_FILE string = "view.html"
)

var (
	WIKI_DIRECTORY     string = DEFAULT_WIKI_DIRECTORY
	CONTENT_DIRECTORY  string = DEFAULT_CONTENT_DIRECTORY
	HTML_TEMPLATE_FILE string = DEFAULT_HTML_TEMPLATE_FILE_FILE
	HTML_TEMPLATE_NAME string = ""
	TEMPLATES          *template.Template
)

type PageNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Children []PageNode `json:"children"`
}

type Page struct {
	Title   string
	Body    template.HTML
	Raw     string
	Sidebar []PageNode
	Session bool
}

func getUrlForPage(directory, filename string) string {
	filename = strings.Replace(filename, ".md", "", -1)
	path := fmt.Sprintf("%v/%v", directory, filename)
	path = strings.Replace(path, WIKI_DIRECTORY, "", -1)
	return path
}

func getDirectoryTree(directory string) []PageNode {
	nodes := []PageNode{}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		logger.Error(err)
		return nodes
	}

	for _, f := range files {
		node := PageNode{Name: f.Name(), Path: getUrlForPage(directory, f.Name())}
		if f.IsDir() {
			node.Children = getDirectoryTree(fmt.Sprintf("%v/%v", directory, f.Name()))
		}
		nodes = append(nodes, node)
	}

	return nodes
}

func buildSideBar() []PageNode {
	return getDirectoryTree(WIKI_DIRECTORY)
}

type WikiEngine struct{}

func (self *WikiEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.ViewHandler(w, r)
}

func (self *WikiEngine) getFilename(page string) string {
	return fmt.Sprintf("%v/%v.md", WIKI_DIRECTORY, page)
}

func (self *WikiEngine) loadPage(page string) (*Page, error) {
	filename := self.getFilename(page)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{
		Title:   page,
		Body:    template.HTML(blackfriday.MarkdownCommon([]byte(body))),
		Raw:     string(body),
		Sidebar: buildSideBar(),
	}, nil
}

func (self *WikiEngine) ViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	page := r.URL.Path[1:]
	if len(page) == 0 {
		page = "index"
	}

	if Sessions.HasSession(r) {
		if "POST" == r.Method {
			err := self.savePage(page, r)
			if err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, string(`{"status":"ok"}`))
			return
		} else if "DELETE" == r.Method {
			err := self.deletePage(page)
			if err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, string(`{"status":"ok"}`))
			return
		}
	}

	p, err := self.loadPage(page)
	if err != nil {
		self.renderTemplate(w, HTML_TEMPLATE_NAME, &Page{Title: page, Sidebar: buildSideBar(), Session: Sessions.HasSession(r)})
		return
	}
	p.Session = Sessions.HasSession(r)
	self.renderTemplate(w, HTML_TEMPLATE_NAME, p)
}

func (self *WikiEngine) deletePage(page string) error {
	// split file path into parts
	path := fmt.Sprintf("%s/%s.md", WIKI_DIRECTORY, page)
	// delete file
	return os.Remove(path)
}

func (self *WikiEngine) savePage(page string, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// split file path into parts
	parts := strings.Split(page, "/")

	// create directory tree from path
	path := fmt.Sprintf("%s/%s", WIKI_DIRECTORY, strings.Join(parts[:len(parts)-1], "/"))
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	// write data to file
	out_file := fmt.Sprintf("%s/%s.md", WIKI_DIRECTORY, strings.Join(parts, "/"))
	err = ioutil.WriteFile(out_file, []byte(body), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (self *WikiEngine) renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := TEMPLATES.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
