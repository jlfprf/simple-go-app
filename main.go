package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

const sessionCookieName = "appck"

var tmpl *template.Template

func main() {

	var e error
	tmpl, e = template.ParseGlob("views/*.html")
	if e != nil {
		panic("Error while parsing templates -> " + e.Error())
	}

	// postgresql://user:secret@localhost:5432/dbname
	// db, err := sql.Open("postgres", "user=postgres password=postgres dbname=simple_go_app sslmode=disable")
	db, err := sql.Open("postgres", "host=172.17.0.2 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not create the database pool with sql.Open().")
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//http.Handle("/static/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler(db))
	http.HandleFunc("/private", privateHandler(db))
	http.HandleFunc("/error", errorHandler)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "This is from a closure")
	})
	http.HandleFunc("/usertest", usertestHandler(db))

	err = http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err.Error())
		panic("Error while trying to listen at port 8080.")
	}
}

//------------------------------------To redirect https--------------------------
// 		h := r.Header
// 		if h.Get("x-forwarded-proto") != "https" {
// 			http.Redirect(w, r, "https://"+r.Host+r.URL.Path, http.StatusFound)
// 			return
// 		}
//-------------------------------------------------------------------------------

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{"Title": "Default Templating with Maps"})
	checkError(err, &w, r)
}

func loginHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := authenticate(w, r, db)
			if err == nil {
				http.Redirect(w, r, "/private", http.StatusFound)
				return
			}
		}
		err := tmpl.ExecuteTemplate(w, "login.html", map[string]string{"Title": "Default Golang Templating - Login"})
		checkError(err, &w, r)
	}
}

func privateHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, ok := isAuthenticated(r, db)
		if !ok {
			// fmt.Println(err.Error())
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		err := tmpl.ExecuteTemplate(w, "private.html", map[string]string{"Title": "Private Area", "User": usr})
		checkError(err, &w, r)
		return
	}

}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "error.html", struct{ Title string }{Title: "Default Go Templating"})
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

//===========================Utils===================================
func checkError(err error, w *http.ResponseWriter, r *http.Request) {
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(*w, r, "/error", http.StatusFound)
		return
	}
}

func authenticate(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	u := r.PostFormValue("u")
	p := r.PostFormValue("p")
	var username, hashedpass string
	_ = db.QueryRow("select name, hashedpass from users where name = $1", u).Scan(&username, &hashedpass)
	if err := bcrypt.CompareHashAndPassword([]byte(hashedpass), []byte(p)); err == nil {
		c := make([]byte, 32)
		_, err := rand.Read(c)
		if err == nil {
			ssessionID := base64.URLEncoding.EncodeToString(c)
			_, err := db.Exec("delete from sessions where name=$1", username)
			if err != nil {
				return errors.New("Error while trying to delete old session from db.")
			}
			_, err = db.Exec("insert into sessions (sessionid, name) values ($1, $2)", ssessionID, username)
			if err == nil {
				ck := &http.Cookie{Name: sessionCookieName, Value: ssessionID, HttpOnly: true}
				http.SetCookie(w, ck)
				// http.Redirect(w, r, "/private", http.StatusFound)
				return nil
			}
		}
	}
	return errors.New("There was an error while authenticating the user.")
}

func isAuthenticated(r *http.Request, db *sql.DB) (string, bool) {
	ck, err := r.Cookie(sessionCookieName)
	if err == nil {
		var result, user string
		row := db.QueryRow("select sessionid, name from sessions where sessionid = $1", ck.Value)
		row.Scan(&result, &user)
		if result != "" {
			return user, true
		}
	}
	return "", false
}

//------------------------------------------Testing--------------------------------------------------

func usertestHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var username, hashedpass string
		_ = db.QueryRow("select name, hashedpass from users where name = $1", "jlf").Scan(&username, &hashedpass)
		fmt.Fprint(w, "User: "+username+"\nP: "+hashedpass)
	}
}
