package server

import (
	h "forum/handlers"
	"net/http"
)

func Routes() *http.ServeMux {
	r := http.NewServeMux()
	r.Handle("/registration", http.HandlerFunc(h.Registration))
	r.Handle("/login", http.HandlerFunc(h.Login))
	r.Handle("/register", http.HandlerFunc(h.Register))
	r.Handle("/", http.HandlerFunc(h.HomePage))
	r.Handle("/logging", http.HandlerFunc(h.Loggin_in))
	r.Handle("/profile", http.HandlerFunc(h.Profile_page))
	r.Handle("/logout", http.HandlerFunc(h.Logout))
	r.Handle("/post_page", http.HandlerFunc(h.Post_page))
	r.Handle("/post", http.HandlerFunc(h.Post))

	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return r
}
