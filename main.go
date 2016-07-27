package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var tmplToParse = []string{"index", "error", "login"}
var tmplsParsed = make(map[string]*template.Template)

func main() {

	tmplsParsed = createTemplates(tmplToParse)

	db, err := sql.Open("postgres", "user=postgres password=mdibhf dbname=pgdatabase")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler(db))
	http.HandleFunc("/db", dbHandler(db))
	http.HandleFunc("/changedata", changeDBData(db))
	http.HandleFunc("/error", errorHandler)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "This is from a closure")
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// c := r.Cookie("appck")
	// if checkAuth(c.Value, db){}
	// http.SetCookie(w, cookie)
	err := tmplsParsed["index"].ExecuteTemplate(w, "Layout", map[string]interface{}{"Title": "Default Templating with Maps"})
	checkError(err, &w, r)
}

func loginHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.PostFormValue("u")
		p := r.PostFormValue("p")
		result := authenticate(u, p, db)
		if result != "" {
			ck := &http.Cookie{Name: "appck", Value: result, HttpOnly: true}
			http.SetCookie(w, ck)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		err := tmplsParsed["login"].ExecuteTemplate(w, "Layout", map[string]string{"Title": "Default Golang Templating"})
		checkError(err, &w, r)
	}
}

func dbHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var name, email string
		_ = db.QueryRow("select name, email from users where name = $1", "jlf").Scan(&name, &email)
		s := fmt.Sprintf("Name: %s\nEmail: %s", name, email)
		fmt.Fprint(w, s)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	err := tmplsParsed["error"].ExecuteTemplate(w, "Layout", struct{ Title string }{Title: "Default Go Templating"})
	checkError(err, &w, r)
}

//===========================Utils===================================
func checkError(err error, w *http.ResponseWriter, r *http.Request) {
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(*w, r, "/error", http.StatusFound)
		return
	}
}

func createTemplates(tmplToParse []string) map[string]*template.Template {
	var tmpls = make(map[string]*template.Template)
	for i := range tmplToParse {
		t, err := template.ParseFiles("views/layout.html", "views/"+tmplToParse[i]+".html")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		tmpls[tmplToParse[i]] = t
	}
	return tmpls
}

func authenticate(u, p string, db *sql.DB) string {
	var username, hashedpass string
	_ = db.QueryRow("select name, hashedpass from users where name = $1", u).Scan(&username, &hashedpass)
	if bcrypt.CompareHashAndPassword([]byte(hashedpass), []byte(p)) == nil {
		c := make([]byte, 32)
		_, err := rand.Read(c)
		if err == nil {
			ckValue := fmt.Sprintf("%x", c)
			_, err := db.Exec("insert into sessions (cookie, name) values ($1, $2)", ckValue, username)
			if err == nil {
				return ckValue
			}
		}
	}
	return ""
	//confirm credentials and save session into db
	//return random text session to create a cookie
}

// func checkAuth(textRandon string, db *sql.DB) bool {
// 	//get data from session cookie and check to see if it is in db
// }

//----------------------------Testing-----------------------------------
func changeDBData(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.FormValue("u")
		e := r.FormValue("e")
		if u != "" && e != "" {
			fmt.Println(u, e)
			_, err := db.Exec("insert into users (name, email) values ($1, $2)", u, e)
			if err != nil {
				fmt.Fprint(w, err.Error())
			}
			fmt.Fprint(w, "Database updated")
			return
		}
		fmt.Fprint(w, "You need to provide a user.")
	}
}
