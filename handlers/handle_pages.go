package handlers

import (
	"net/http"
	"strconv"
	"text/template"
)

type User struct {
	User_email interface{}
	User_name  interface{}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := Store.Get(r, "user")
	username := session.Values["User_name"]
	email, ok := session.Values["User_email"]
	var user User
	if ok {
		user.User_name = username
		user.User_email = email
	}

	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	tmpl, err := template.ParseFiles("./static/templates/home.html")
	if err != nil {
		ErrorLog.Print("Error parsing files\n", err.Error)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, user)
	if err != nil {
		ErrorLog.Print("Error template exexuting\n", err.Error)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	tmpl, err := template.ParseFiles("./static/templates/login.html")
	if err != nil {
		ErrorLog.Print("Error parsing files\n", err.Error)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, ""); err != nil {
		ErrorLog.Print("Error template exexuting\n", err.Error)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func Registration(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/registration" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	tmpl, err := template.ParseFiles("./static/templates/reg-page.html")
	if err != nil {
		ErrorLog.Print("Error parsing files\n", err.Error)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, ""); err != nil {
		ErrorLog.Print("Error template exexuting\n", err.Error)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	tmpl, err := template.ParseFiles("./static/templates/error.html")
	if err != nil {
		w.Write([]byte(http.StatusText(status)))
		return
	}
	tmpl.Execute(w, http.StatusText(status)+" "+strconv.Itoa(status))
}

func Profile_page(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := Store.Get(r, "user")
	username := session.Values["User_name"]
	email, ok := session.Values["User_email"]
	var user User
	if ok {
		user.User_name = username
		user.User_email = email
	}

	tmpl, err := template.ParseFiles("./static/templates/profile_page.html")
	if err != nil {
		ErrorLog.Print("Error parsing files\n", err.Error)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, user); err != nil {
		ErrorLog.Print("Error template executing\n", err.Error)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "user")
	delete(session.Values, "User_email")
	delete(session.Values, "User_name")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
