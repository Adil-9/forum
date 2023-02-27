package handlers

import (
	"net/http"
	"strconv"
	"text/template"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
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
	if err := tmpl.Execute(w, ""); err != nil {
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
	tmpl, err := template.ParseFiles("./static/templates/login.htm")
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

func Internal_error(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./static/templates/error.html")
	tmpl.Execute(w, http.StatusText(404)+" 404")
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	tmpl, _ := template.ParseFiles("./static/templates/error.html")
	tmpl.Execute(w, http.StatusText(status)+" "+strconv.Itoa(status))
}
