package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var tmplToParse = []string{"index", "error", "login", "private"}
var tmplsParsed = make(map[string]*template.Template)

const sessionCookieName = "appck"

func main() {

	tmplsParsed = createTemplates(tmplToParse)

	db, err := sql.Open("postgres", "user=postgres password=mdibhf dbname=simple_go_app")
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not create the database pool with sql.Open().")
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler(db))
	http.HandleFunc("/private", privateHandler(db))
	http.HandleFunc("/error", errorHandler)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "This is from a closure")
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
		// os.Exit(1)
		panic("Error while trying to listen at port 8080.")
	}
}

//------------------------------------To redirect https--------------------------
// if IN_HEROKU {
// 		h := r.Header
// 		if h.Get("x-forwarded-proto") != "https" {
// 			http.Redirect(w, r, "https://"+r.Host+r.URL.Path, http.StatusFound)
// 			return
// 		}
// 	}
//-------------------------------------------------------------------------------

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := tmplsParsed["index"].ExecuteTemplate(w, "Layout", map[string]interface{}{"Title": "Default Templating with Maps"})
	checkError(err, &w, r)
}

func loginHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && authenticate(w, r, db) {
			return
		}
		err := tmplsParsed["login"].ExecuteTemplate(w, "Layout", map[string]string{"Title": "Default Golang Templating - Login"})
		checkError(err, &w, r)
	}
}

func privateHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if usr, ok := isAuthenticated(r, db); ok {
			err := tmplsParsed["private"].ExecuteTemplate(w, "Layout", map[string]string{"Title": "Private Area", "User": usr})
			checkError(err, &w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
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
			panic("Could not process the templates. Func createTemplates()")
		}
		tmpls[tmplToParse[i]] = t
	}
	return tmpls
}

func authenticate(w http.ResponseWriter, r *http.Request, db *sql.DB) bool {
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
				return false
			}
			_, err = db.Exec("insert into sessions (sessionid, name) values ($1, $2)", ssessionID, username)
			if err == nil {
				ck := &http.Cookie{Name: sessionCookieName, Value: ssessionID, HttpOnly: true}
				http.SetCookie(w, ck)
				http.Redirect(w, r, "/private", http.StatusFound)
				return true
			}
		}
	}
	return false
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
