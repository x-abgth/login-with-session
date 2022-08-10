package main

import (
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	tpl *template.Template

	key = []byte("brototype-first-project")
	// NewCookieStore returns a new cookie store
	// CookieStore stores sessions using secure cookies.
	store = sessions.NewCookieStore(key)
)

func init() {
	var err error

	// Parsing every .html files in the folder ../views/
	tpl, err = tpl.ParseGlob("./views/*.html")

	if err != nil {
		panic(err)
	}
}

func main() {
	router := mux.NewRouter()

	// This code specifies the assets path - ../views/assets/
	fs := http.FileServer(http.Dir("./views/assets/"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// accessing handlers
	router.HandleFunc("/", signInPageHandler)
	router.HandleFunc("/authenticate", authenticateHandler)
	router.HandleFunc("/home", homeHandler)
	router.HandleFunc("/signout", signOutPageHandler)
	router.Handle("/favicon.ico", router.NotFoundHandler)

	// Listen and serve will accept two values - port number and a mux
	log.Fatal(http.ListenAndServe(":8080", router))
}
