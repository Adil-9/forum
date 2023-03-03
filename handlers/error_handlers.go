package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

var (
	InfoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	tmpl, err := template.ParseFiles("./static/templates/error.html")
	if err != nil {
		w.Write([]byte(http.StatusText(status)))
		return
	}
	tmpl.Execute(w, http.StatusText(status)+" "+strconv.Itoa(status))
}
