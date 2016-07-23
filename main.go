package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var tmpls = make(map[string]*template.Template)

func main() {

	tmpls = createTemplates()

	db, err := sql.Open("postgres", "user=postgres password=mdibhf dbname=pgdatabase")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/db", dbHandler(db))
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
	err := tmpls["index"].ExecuteTemplate(w, "Base", map[string]interface{}{"Title": "Default Templating with Maps"})
	checkError(err, &w, r)
}

func dbHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var name, email string
		_ = db.QueryRow("select name, email from users where name = $1", "rctf").Scan(&name, &email)
		s := fmt.Sprintf("Name: %s\nEmail: %s", name, email)
		fmt.Fprint(w, s)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpls["error"].ExecuteTemplate(w, "Base", struct{ Title string }{Title: "Default Go Templating"})
	checkError(err, &w, r)
}

func checkError(err error, w *http.ResponseWriter, r *http.Request) {
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(*w, r, "/error", http.StatusFound)
		return
	}
}

func createTemplates() map[string]*template.Template {
	var tmpls = make(map[string]*template.Template)
	t, err := template.ParseFiles("views/base.html", "views/index.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	tmpls["index"] = t
	t, err = template.ParseFiles("views/base.html", "views/error.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	tmpls["error"] = t
	return tmpls
}
