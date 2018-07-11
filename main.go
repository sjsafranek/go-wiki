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

var users_file string

func init() {
	flag.IntVar(&PORT, "p", DEFAULT_PORT, "Server port")
	flag.StringVar(&CONTENT_DIRECTORY, "C", DEFAULT_CONTENT_DIRECTORY, "Wiki content directory")
	flag.StringVar(&HTML_TEMPLATE_FILE, "t", DEFAULT_HTML_TEMPLATE_FILE_FILE, "Html template")
	// flag.StringVar(p, name, value, usage)
	flag.Parse()

	TEMPLATES = template.Must(template.ParseFiles(HTML_TEMPLATE_FILE))
	HTML_TEMPLATE_NAME = strings.Replace(HTML_TEMPLATE_FILE, ".html", "", -1)
}

func main() {
	var err error

	go RunTcpServer()

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

	// Sessions
	Sessions = AuthenticationMiddleware{File: fmt.Sprintf("%v/users.json", CONTENT_DIRECTORY)}
	Sessions.Init()
	router.HandleFunc("/login", Sessions.LoginHandler).Methods("GET", "POST")
	router.HandleFunc("/logout", Sessions.LogoutHandler).Methods("GET")
	//.end

	// File uploader
	router.Handle("/upload", Sessions.RequireSession(http.HandlerFunc(FileUploadHandler))).Methods("GET", "POST")
	//.end

	// wiki engine
	wiki := &WikiEngine{}
	WIKI_DIRECTORY = fmt.Sprintf("%v/wiki/", CONTENT_DIRECTORY)
	err = os.MkdirAll(WIKI_DIRECTORY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	router.PathPrefix("/").Handler(wiki)
	//.end

	router.Use(LoggingMiddleWare, SetHeadersMiddleWare, Sessions.Middleware)

	// http.Handle("/static/", http.FileServer(http.Dir("static")))
	// fs_static := http.FileServer(http.Dir("static"))
	// static_route := fmt.Sprintf("/%v/", "static")
	// http.Handle(static_route, http.StripPrefix(static_route, fs_static))

	logger.Infof("Magic happens on port %v...", PORT)
	// err = http.ListenAndServe(fmt.Sprintf(":%v", PORT), TrafficMiddleWare(router))
	err = http.ListenAndServe(fmt.Sprintf(":%v", PORT), router)
	if nil != err {
		panic(err)
	}
}
