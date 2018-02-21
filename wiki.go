package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/russross/blackfriday"
)

const (
	DEFAULT_CONTENT_DIRECTORY       string = "content"
	DEFAULT_PORT                    int    = 1337
	DEFAULT_HTML_TEMPLATE_FILE_FILE string = "view.html"
)

var (
	CONTENT_DIRECTORY  string = DEFAULT_CONTENT_DIRECTORY
	PORT               int    = DEFAULT_PORT
	HTML_TEMPLATE_FILE string = DEFAULT_HTML_TEMPLATE_FILE_FILE
	HTML_TEMPLATE_NAME string = ""
	TEMPLATES          *template.Template
	USERS              Users
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
		err := errors.New("Page not found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, err := self.loadPage(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

//
// func init() {
// 	flag.IntVar(&PORT, "p", DEFAULT_PORT, "Server port")
// 	flag.StringVar(&CONTENT_DIRECTORY, "C", DEFAULT_CONTENT_DIRECTORY, "Content directory")
// 	flag.StringVar(&HTML_TEMPLATE_FILE, "t", DEFAULT_HTML_TEMPLATE_FILE_FILE, "Html template")
// 	flag.Parse()
//
// 	TEMPLATES = template.Must(template.ParseFiles(HTML_TEMPLATE_FILE))
// 	HTML_TEMPLATE_NAME = strings.Replace(HTML_TEMPLATE_FILE, ".html", "", -1)
//
// 	// set wiki users
// 	USERS = Users{}
// 	USERS.Fetch("users.json")
// 	user := User{Username: "admin"}
// 	USERS.Add(&user)
// 	user.SetPassword("dev")
// 	USERS.Save("users.json")
// 	//.end
//
// }
//
// func main() {
// 	// http://www.alexedwards.net/blog/a-recap-of-request-handling
//
// 	wiki := &WikiEngine{}
// 	// http.Handle("/wiki/", http.StripPrefix("/wiki", wiki))
// 	http.Handle("/", wiki)
//
// 	// Static Files
// 	fs := http.FileServer(http.Dir("img"))
// 	http.Handle("/img/", http.StripPrefix("/img/", fs))
//
// 	fmt.Printf("Magic happens on port %v...\n", PORT)
// 	err := http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
// 	if nil != err {
// 		panic(err)
// 	}
// }
