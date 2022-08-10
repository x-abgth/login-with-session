package main

import (
	"net/http"
	"time"
)

func signInPageHandler(w http.ResponseWriter, r *http.Request) {
	// Getting session
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Checking the authenticated key's value in the session
	if session.Values["authenticated"] == true && session.Values["isInvalid"] == false {
		// If authenticated, then redirect to home page
		http.Redirect(w, r, "/home", http.StatusFound)
	} else if session.Values["isInvalid"] == true {
		// if the email and password is wrong, show error message!
		val := credentials{Title: "Login Page", ErrMsg: "Invalid email id or password"}
		tpl.ExecuteTemplate(w, "sign_in.html", val)
	} else {
		// First clear cache
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.

		// Then show the login page,
		val := credentials{Title: "Login Page"}
		tpl.ExecuteTemplate(w, "sign_in.html", val)
	}
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	// Parsing form to get the values we entered
	r.ParseForm()

	// response object r will give us the entered value
	email := r.PostForm.Get("emailId")
	pass := r.PostForm.Get("passwordVal")

	// Getting session named 'session'
	session, _ := store.Get(r, "session")

	// validating email and password
	if email == "test@gmail.com" && pass == "12345" {

		// Assigning values to the session variables authenticated and emailId
		session.Values["authenticated"] = true
		session.Values["isInvalid"] = false
		session.Values["emailId"] = email
		// Saving the session values
		session.Save(r, w)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		// if password and username is wrong, we assign false to authenticated
		session.Values["authenticated"] = false
		session.Values["isInvalid"] = true
		// saving the session variables
		session.Save(r, w)

		// Redirect to the login page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Get session
	session, _ := store.Get(r, "session")

	// If the user is already authenticated, then show home page
	if session.Values["authenticated"] == true {
		email := session.Values["emailId"].(string)

		val := credentials{
			Title: "Home Page",
			Email: email,
		}

		// Clear cache
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.
		// Then get to home page
		err := tpl.ExecuteTemplate(w, "home.html", val)
		if err != nil {
			panic(err)
		}
	} else {
		// Clear the cache
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.
		// Redirect to login page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func signOutPageHandler(w http.ResponseWriter, r *http.Request) {
	// Clearing the cache memory
	w.Header().Set("Cache-Control", "no-cache, private, no-store, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")

	// Getting the session
	session, _ := store.Get(r, "session")

	// Destroying the sessions by assigning the maxAge value to -1
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1

	// Saving the session
	session.Save(r, w)

	// Redirect to the login page
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
