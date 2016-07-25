package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

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
	http.HandleFunc("/login", loginHandler)
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
	err := tmplsParsed["index"].ExecuteTemplate(w, "Base", map[string]interface{}{"Title": "Default Templating with Maps"})
	checkError(err, &w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := tmplsParsed["login"].ExecuteTemplate(w, "Base", map[string]string{"Title": "Default Golang Templating"})
	checkError(err, &w, r)
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
	err := tmplsParsed["error"].ExecuteTemplate(w, "Base", struct{ Title string }{Title: "Default Go Templating"})
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
		t, err := template.ParseFiles("views/base.html", "views/"+tmplToParse[i]+".html")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		tmpls[tmplToParse[i]] = t
	}
	return tmpls
}

// func authenticate(u, p string, db *sql.DB) string {
// 	//confirm credentials and save session into db
// 	//return random text session to create a cookie
// }

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
