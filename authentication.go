package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type AuthenticationMiddleware struct {
	store *sessions.CookieStore
	File  string
	db    Users
}

func (self *AuthenticationMiddleware) HasSession(r *http.Request) bool {
	session, err := self.store.Get(r, "wiki-session")
	if nil != err {
		return false
	}
	if nil == session.Values["loggedin"] {
		return false
	}
	return true
}

func (self *AuthenticationMiddleware) CreateUser(username, password string) {
	user := User{Username: username}
	user.SetPassword(password)
	self.db.Add(&user)
	self.db.Save(self.File)
}

func (self *AuthenticationMiddleware) Init() {
	// generate new secret for sessions every time the server starts
	secret := uuid.New().String()
	self.store = sessions.NewCookieStore([]byte(secret))

	// set wiki users
	err := os.MkdirAll(CONTENT_DIRECTORY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	self.db = Users{}
	self.db.Fetch(self.File)
	user := User{Username: "admin"}
	user.SetPassword("dev")
	self.db.Add(&user)
	self.db.Save(self.File)
	//.end
}

// Middleware function, which will be called for specified requests
func (self *AuthenticationMiddleware) RequireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !self.HasSession(r) {
			logger.Warn("Not authenticated")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware function, which will be called for each request
func (self *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "POST" == r.Method && "/login" != r.URL.Path {

			session, err := self.store.Get(r, "wiki-session")
			if nil != err {
				logger.Warn("Not authenticated")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			if nil == session.Values["loggedin"] {
				logger.Warn("Not authenticated")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

		}

		next.ServeHTTP(w, r)
	})
}

func (self *AuthenticationMiddleware) loginTemplate(err error) string {
	err_msg := ""
	if nil != err {
		err_msg = err.Error()
	}
	return `<!DOCTYPE html>
				<html>
					<head></head>
					<body>
						<h3>Login</h3>
						<form action="/login" method="POST">
							Username: <input type="text" placeholder="username" name="username"><br>
							Password: <input type="password" name="password"><br>
							<input type="submit" value="Login">
							<div>` + err_msg + `</div>
						</form>
					<body>
				</html>`
}

func (self *AuthenticationMiddleware) LoginHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {
		fmt.Fprintf(w, self.loginTemplate(nil))
		return
	}

	err := r.ParseForm()
	if nil != err {
		http.Error(w, "Unable to parse form", http.StatusInternalServerError)
		return
	}

	username := r.Form["username"][0]
	user, err := self.db.Get(username)
	if nil != err {
		logger.Warn(err)
		fmt.Fprintf(w, self.loginTemplate(err))
		return
	}

	password := r.Form["password"][0]
	if !user.IsPassword(password) {
		err = errors.New("Incorrect password")
		logger.Warn(err)
		fmt.Fprintf(w, self.loginTemplate(err))
		return
	}

	// create session
	session, _ := self.store.Get(r, "wiki-session")
	session.Values["loggedin"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (self *AuthenticationMiddleware) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// delete session
	session, _ := self.store.Get(r, "wiki-session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
