package main

import (
	"flag"
	h "forum/handlers"
	s "forum/server"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	srv := http.Server{
		Addr:     *addr,
		ErrorLog: h.ErrorLog,
		Handler:  s.Routes(),
	}

	h.InfoLog.Printf("Server running on: http://localhost%s", *addr)
	err := srv.ListenAndServe()
	h.ErrorLog.Fatal(err)
}
