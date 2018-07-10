package main

import (
	// "errors"
	"fmt"
	"html/template"
	"io/ioutil"
	// "log"
	"net/http"
	"strings"

	"os"
	// "path/filepath"

	"github.com/russross/blackfriday"

	"github.com/sjsafranek/ligneous"
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
	logger             = ligneous.New().Log
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
}

func getUrlForPage(directory, filename string) string {
	filename = strings.Replace(filename, ".md", "", -1)
	path := fmt.Sprintf("%v/%v", directory, filename)
	path = strings.Replace(path, CONTENT_DIRECTORY, "", -1)
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
	// tree := []Node{}
	// err := filepath.Walk(CONTENT_DIRECTORY, func(path string, f os.FileInfo, err error) error {
	// 	fmt.Printf("Visited: %s %s %s\n", path, f.IsDir(), f.Name())
	//
	// 	return nil
	// })
	// if nil != err {
	// 	fmt.Println(err)
	// }

	return getDirectoryTree(CONTENT_DIRECTORY)
}

type WikiEngine struct{}

func (self *WikiEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.URL.Path)
	logger.Infof("%s [%s]", r.URL.Path, r.Method)
	// TODO:
	// - add /edit route,
	// - POST request to create new pages
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

	return &Page{
		Title:   page,
		Body:    template.HTML(blackfriday.MarkdownCommon([]byte(body))),
		Raw:     string(body),
		Sidebar: buildSideBar(),
	}, nil
}

func (self *WikiEngine) viewHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path[1:]
	if len(page) == 0 {
		page = "index"
	}

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
	}

	p, err := self.loadPage(page)
	if err != nil {
		self.renderTemplate(w, HTML_TEMPLATE_NAME, &Page{})
		return
	}
	self.renderTemplate(w, HTML_TEMPLATE_NAME, p)
}

func (self *WikiEngine) savePage(page string, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// split file path into parts
	parts := strings.Split(page, "/")

	// create directory tree from path
	path := fmt.Sprintf("%s/%s", CONTENT_DIRECTORY, strings.Join(parts[:len(parts)-1], "/"))
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	// write data to file
	out_file := fmt.Sprintf("%s/%s.md", CONTENT_DIRECTORY, strings.Join(parts, "/"))
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
