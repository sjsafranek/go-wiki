package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func HasSession(r *http.Request) bool {
	session, err := store.Get(r, "wiki-session")
	if nil != err {
		return false
	}
	if nil == session.Values["loggedin"] {
		return false
	}
	return true
}

func init() {
	// generate new secret for sessions every time the server starts
	secret := uuid.New().String()
	store = sessions.NewCookieStore([]byte(secret))
}

type middleware struct {
	handler http.Handler
}

func (self middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// set default headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	self.handler.ServeHTTP(w, r)
	logger.Debug(fmt.Sprintf("%v %v %v", r.RemoteAddr, r.Method, r.URL))
}

func MiddleWare(h http.Handler) http.Handler {
	return middleware{h}
}

type AuthenticationMiddleware struct {
	tokens map[string]string
}

func (self *AuthenticationMiddleware) Populate() {
	self.tokens = make(map[string]string)
	for _, user := range USERS.Users {
		self.tokens[user.Password] = user.Username
	}
}

// Middleware function, which will be called for each request
func (self *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// token := r.Header.Get("X-Session-Token
		logger.Info(r.URL.Path[1:], r.URL.Path)

		if "POST" == r.Method && "/login" != r.URL.Path {
			// if user, found := self.tokens[token]; found {
			// 	logger.Infof("Authenticated user %s\n", user)
			// 	next.ServeHTTP(w, r)
			// } else {
			// 	logger.Warnf("Not authenticated user %s\n", user)
			// 	http.Error(w, "Forbidden", http.StatusForbidden)
			// }
			// return
			session, err := store.Get(r, "wiki-session")
			if nil != err {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			if nil == session.Values["loggedin"] {
				logger.Warn("Not authenticated")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			logger.Info("Authenticated")
		}

		next.ServeHTTP(w, r)
	})
}

func (self *AuthenticationMiddleware) LoginHandler(w http.ResponseWriter, r *http.Request) {

	if "GET" == r.Method {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
				<head></head>
				<body>
					<form action="/login" method="POST">
						Username: <input type="text" placeholder="username" name="username"><br>
						Password: <input type="password" name="password"><br>
						<input type="submit" value="Login">
					</form>
				<body>
			</html>
		`)
		return
	}

	err := r.ParseForm()
	if nil != err {
		http.Error(w, "Unable to parse form", http.StatusInternalServerError)
		return
	}

	username := r.Form["username"][0]
	password := r.Form["password"][0]
	user, err := USERS.Get(username)
	if nil != err {
		// http.Error(w, err.Error(), http.StatusNotFound)
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
				<head></head>
				<body>
					<form action="/login" method="POST">
						Username: <input type="text" placeholder="username" name="username"><br>
						Password: <input type="password" name="password"><br>
						<input type="submit" value="Login">
						<div>`+err.Error()+`</div>
					</form>
				<body>
			</html>
		`)
		return
	}

	if !user.IsPassword(password) {
		// http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
				<head></head>
				<body>
					<form action="/login" method="POST">
						Username: <input type="text" placeholder="username" name="username"><br>
						Password: <input type="password" name="password"><br>
						<input type="submit" value="Login">
						<div>Forbidden</div>
					</form>
				<body>
			</html>
		`)
		return
	}

	// create session
	session, _ := store.Get(r, "wiki-session")
	session.Values["loggedin"] = true
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(`{"status":"ok"}`))
}

func (self *AuthenticationMiddleware) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// delete session
	session, _ := store.Get(r, "wiki-session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	// w.Header().Set("Content-Type", "application/json")
	// fmt.Fprintf(w, string(`{"status":"ok"}`))
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)

}
