package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func init() {
	flag.IntVar(&PORT, "p", DEFAULT_PORT, "Server port")
	flag.StringVar(&CONTENT_DIRECTORY, "C", DEFAULT_CONTENT_DIRECTORY, "Wiki content directory")
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

	router := mux.NewRouter()

	// http://www.alexedwards.net/blog/a-recap-of-request-handling

	// TODO
	//  - Create util function for this
	// Static Files
	IMG_DIRECTORY = fmt.Sprintf("%v/img/", CONTENT_DIRECTORY)
	err = os.MkdirAll(IMG_DIRECTORY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	router.PathPrefix("/img/").Handler(
		http.StripPrefix("/img/", http.FileServer(
			http.Dir(IMG_DIRECTORY))))
	//.end

	// File uploader
	router.HandleFunc("/upload", FileUploadHandler).Methods("GET", "POST")
	//.end

	auth := AuthenticationMiddleware{}
	auth.Populate()
	router.HandleFunc("/login", auth.LoginHandler).Methods("GET", "POST")
	router.HandleFunc("/logout", auth.LogoutHandler).Methods("GET")

	// wiki engine
	wiki := &WikiEngine{}
	WIKI_DIRECTORY = fmt.Sprintf("%v/wiki/", CONTENT_DIRECTORY)
	err = os.MkdirAll(WIKI_DIRECTORY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	router.PathPrefix("/").Handler(wiki)
	//.end

	router.Use(auth.Middleware)

	// http.Handle("/static/", http.FileServer(http.Dir("static")))
	// fs_static := http.FileServer(http.Dir("static"))
	// static_route := fmt.Sprintf("/%v/", "static")
	// http.Handle(static_route, http.StripPrefix(static_route, fs_static))

	fmt.Printf("Magic happens on port %v...\n", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%v", PORT), MiddleWare(router))
	if nil != err {
		panic(err)
	}
}
