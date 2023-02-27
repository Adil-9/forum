package server

import (
	"flag"
	h "forum/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func Server() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	r := mux.NewRouter()
	r.Handle("/registration", http.HandlerFunc(h.Registration))
	r.Handle("/login", http.HandlerFunc(h.Login))
	r.Handle("/register", http.HandlerFunc(h.Register))
	r.Handle("/", http.HandlerFunc(h.HomePage))
	r.Handle("/logging", http.HandlerFunc(h.Loggin_in))
	r.NotFoundHandler = http.HandlerFunc(h.Internal_error)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	srv := http.Server{
		Addr:     *addr,
		ErrorLog: h.ErrorLog,
		Handler:  r,
	}

	h.InfoLog.Printf("Server running on: http://localhost%s", *addr)
	err := srv.ListenAndServe()
	h.ErrorLog.Println(err.Error())
}
