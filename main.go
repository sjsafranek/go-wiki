package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
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
	// flag.StringVar(&IMG_DIRECTORY, "i", DEFAULT_IMG_DIRECTORY, "Image directory")
	flag.StringVar(&HTML_TEMPLATE_FILE, "t", DEFAULT_HTML_TEMPLATE_FILE_FILE, "Html template")
	// flag.StringVar(p, name, value, usage)
	flag.Parse()

	TEMPLATES = template.Must(template.ParseFiles(HTML_TEMPLATE_FILE))
	HTML_TEMPLATE_NAME = strings.Replace(HTML_TEMPLATE_FILE, ".html", "", -1)

	// set wiki users
	err := os.MkdirAll(CONTENT_DIRECTORY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	users_file := fmt.Sprintf("%v/users.json", CONTENT_DIRECTORY)
	USERS = Users{}
	USERS.Fetch(users_file)
	user := User{Username: "admin"}
	USERS.Add(&user)
	user.SetPassword("dev")
	USERS.Save(users_file)
	//.end
}

func main() {
	var err error

	// http://www.alexedwards.net/blog/a-recap-of-request-handling

	wiki := &WikiEngine{}
	// http.Handle("/wiki/", http.StripPrefix("/wiki", wiki))

	WIKI_DIRECTORY = fmt.Sprintf("%v/wiki/", CONTENT_DIRECTORY)
	err = os.MkdirAll(WIKI_DIRECTORY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	http.Handle("/", wiki)

	// Static Files
	IMG_DIRECTORY = fmt.Sprintf("%v/img/", CONTENT_DIRECTORY)
	err = os.MkdirAll(IMG_DIRECTORY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	fs_image := http.FileServer(http.Dir(IMG_DIRECTORY))
	http.Handle("/img/", http.StripPrefix("/img/", fs_image))

	// http.Handle("/static/", http.FileServer(http.Dir("static")))
	// fs_static := http.FileServer(http.Dir("static"))
	// static_route := fmt.Sprintf("/%v/", "static")
	// http.Handle(static_route, http.StripPrefix(static_route, fs_static))

	fmt.Printf("Magic happens on port %v...\n", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
	if nil != err {
		panic(err)
	}
}
