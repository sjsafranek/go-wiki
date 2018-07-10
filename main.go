package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

const (
	DEFAULT_PORT          int    = 1337
	DEFAULT_IMG_DIRECTORY string = "img"
)

var (
	PORT          int = DEFAULT_PORT
	USERS         Users
	IMG_DIRECTORY string = DEFAULT_IMG_DIRECTORY
)

func init() {
	flag.IntVar(&PORT, "p", DEFAULT_PORT, "Server port")
	flag.StringVar(&CONTENT_DIRECTORY, "C", DEFAULT_CONTENT_DIRECTORY, "Content directory")
	flag.StringVar(&IMG_DIRECTORY, "i", DEFAULT_IMG_DIRECTORY, "Image directory")
	flag.StringVar(&HTML_TEMPLATE_FILE, "t", DEFAULT_HTML_TEMPLATE_FILE_FILE, "Html template")
	// flag.StringVar(p, name, value, usage)
	flag.Parse()

	TEMPLATES = template.Must(template.ParseFiles(HTML_TEMPLATE_FILE))
	HTML_TEMPLATE_NAME = strings.Replace(HTML_TEMPLATE_FILE, ".html", "", -1)

	// set wiki users
	USERS = Users{}
	USERS.Fetch("users.json")
	user := User{Username: "admin"}
	USERS.Add(&user)
	user.SetPassword("dev")
	USERS.Save("users.json")
	//.end
}

func main() {
	// http://www.alexedwards.net/blog/a-recap-of-request-handling

	wiki := &WikiEngine{}
	// http.Handle("/wiki/", http.StripPrefix("/wiki", wiki))

	http.Handle("/", wiki)

	// Static Files
	fs := http.FileServer(http.Dir(IMG_DIRECTORY))
	img_route := fmt.Sprintf("/%v/", IMG_DIRECTORY)
	http.Handle(img_route, http.StripPrefix(img_route, fs))

	// http.Handle("/static/", http.FileServer(http.Dir("static")))
	fs_static := http.FileServer(http.Dir("static"))
	static_route := fmt.Sprintf("/%v/", "static")
	http.Handle(static_route, http.StripPrefix(static_route, fs_static))

	fmt.Printf("Magic happens on port %v...\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
	if nil != err {
		panic(err)
	}
}
