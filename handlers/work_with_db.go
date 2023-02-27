package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("forum_session_key")))

func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	email := r.Form.Get("email-address")
	passwd_1 := r.Form.Get("password-first")
	passwd_2 := r.Form.Get("password-second")
	if passwd_1 != passwd_2 {
		http.Redirect(w, r, "/registration", http.StatusSeeOther)
	}
	db, err := sql.Open("sqlite3", "./databases/logins.db")
	if err != nil {
		fmt.Println("Error opening logs database")
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY,
		username TEXT,
		email TEXT,
		password TEXT
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Println("error sqlstmt")
		panic(err)
	}

	stmt, err := db.Prepare("SELECT email, username FROM logs WHERE email=? OR username=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	exists := stmt.QueryRow(email, username)
	var check_email string
	var check_username string
	exists.Scan(&check_email, &check_username)
	if check_email == email || check_username == username {
		http.Redirect(w, r, "/registration", http.StatusSeeOther)
		return
	}
	stmt, err = db.Prepare("INSERT INTO logs (username, email, password) VALUES(?,?,?)")
	if err != nil {
		fmt.Println("Error prepareing database")
		panic(err)
	}
	defer stmt.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(passwd_1), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("hashing password went wrong")
	}

	_, err = stmt.Exec(username, email, hash) // here parse from html and add
	if err != nil {
		fmt.Println("Error inserting values to database")
		panic(err)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther) // change to main page and give session
}

func Loggin_in(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email-address")
	passwd := r.Form.Get("password")

	db, err := sql.Open("sqlite3", "./databases/logins.db")
	if err != nil {
		panic(err.Error())
	}
	stmt, err := db.Prepare("SELECT password, username FROM logs WHERE email=? LIMIT 1")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(email)
	if err != nil {
		panic(err.Error())
	}
	var password string
	var username string
	for rows.Next() {
		rows.Scan(&password, &username)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(passwd)); err != nil {
		fmt.Println("The password is incorrect")
		http.Redirect(w, r, "/registration", http.StatusSeeOther)
	} else {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, _ := store.Get(r, username)
		// Set some session values.
		session.Values["Username"] = username
		session.Values["Email"] = email
		// Save it before we write to the response/return from the handler.
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
