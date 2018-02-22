package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

const (
	DEFAULT_PORT int = 1337
)

var (
	PORT  int = DEFAULT_PORT
	USERS Users
)

func init() {
	flag.IntVar(&PORT, "p", DEFAULT_PORT, "Server port")
	flag.StringVar(&CONTENT_DIRECTORY, "C", DEFAULT_CONTENT_DIRECTORY, "Content directory")
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

	// USERS.LoadGzip()
	// USERS.SaveGzip()

}

func main() {
	// http://www.alexedwards.net/blog/a-recap-of-request-handling

	wiki := &WikiEngine{}
	// http.Handle("/wiki/", http.StripPrefix("/wiki", wiki))
	http.Handle("/", wiki)

	// Static Files
	fs := http.FileServer(http.Dir("img"))
	http.Handle("/img/", http.StripPrefix("/img/", fs))

	fmt.Printf("Magic happens on port %v...\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
	if nil != err {
		panic(err)
	}
}
