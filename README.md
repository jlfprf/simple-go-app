# simple-go-app
Testing web programming with golang.
My thinking is to use this app as boilerplate for future projects.

##Authentication
In this app I will be using a db table to hold sessions information.
The _**authenticate**_ function compares the data from login form with the user in db, 
then creates a randon string to be saved into session table on db and to used as a session cookie if
authentication fails return a void string.
The _**checkAuth**_ function gets data from cookie session and check to see if it is present in session table on db.

##Templating
It is used the default golang templating. It was created the global string array var _**tmplToParse**_ and 
global template.Template array _**tmplsParsed**_. The tmplToParse array holds the name of files who contains 
templates to be parsed and combined with the layout that is the base. The function _**createTemplates**_ 
creates the templates into    _**tmplsParsed**_ to be used by the handlers. 
```
fmt.Println("this for code")
```

##checkError
This function aims at redirecting the user to a error page when an error happens that can be handled without closing the app altogether.
