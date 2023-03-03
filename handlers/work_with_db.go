package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var Store = sessions.NewCookieStore([]byte(os.Getenv("forum_session_key")))

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
		ErrorLog.Println("Error opening logs database")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// make sure table exists
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY,
		username TEXT,
		email TEXT,
		password TEXT
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		ErrorLog.Println("error sqlstmt")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("SELECT email, username FROM logs WHERE email=? OR username=?")
	if err != nil {
		ErrorLog.Println(err.Error())
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
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
		ErrorLog.Println("Error prepareing database")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(passwd_1), bcrypt.DefaultCost)
	if err != nil {
		ErrorLog.Println("hashing password went wrong")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(username, email, hash) // here parse from html and add
	if err != nil {
		ErrorLog.Println("Error inserting values to database")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := Store.Get(r, "user")
	// Set some session values.
	session.Values["User_name"] = username
	session.Values["User_email"] = email
	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func Loggin_in(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email-address")
	passwd := r.Form.Get("password")

	db, err := sql.Open("sqlite3", "./databases/logins.db")
	if err != nil {
		ErrorLog.Println(err.Error())
	}

	// make sure table exists
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY,
		username TEXT,
		email TEXT,
		password TEXT
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		ErrorLog.Println("error sqlstmt")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("SELECT password, username FROM logs WHERE email=? LIMIT 1")
	if err != nil {
		ErrorLog.Println(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(email)
	if err != nil {
		ErrorLog.Println(err.Error())
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
		session, _ := Store.Get(r, "user")
		// Set some session values.
		session.Values["User_name"] = username
		session.Values["User_email"] = email
		// Save it before we write to the response/return from the handler.
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}

func Post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session, _ := Store.Get(r, "user")
	username, ok := session.Values["User_name"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	date := time.DateTime
	content := r.Form.Get("content")
	db, err := sql.Open("sqlite3", "./databases/posts.db")
	if err != nil {
		ErrorLog.Println("Error opening logs database")
		return
	}
	defer db.Close()

	// make sure table exists
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY,
		username TEXT,
		content TEXT,
		date TEXT,
		likes INTEGER,
		dislikes INTEGER
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		ErrorLog.Println("error sqlstmt")
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("INSERT INTO post (content, username, date) values (?,?,?)")
	if err != nil {
		ErrorLog.Println(err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(content, username, date)
	if err != nil {
		ErrorLog.Println("error sqlstmt")
		return
	}
	http.Redirect(w, r, "/post_page", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "user")
	delete(session.Values, "User_email")
	delete(session.Values, "User_name")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
