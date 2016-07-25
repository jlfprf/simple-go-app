# simple-go-app
Testing web programming with golang.

#Authentication
It is necessary to create a session table on the db.
The _**authenticate**_ function compares the data from login form with the user in db, then creates a randon string to be saved into session table on db and to used as a session cookie.
The checkAuth function gets data from cookie session and check to see if it is present in session table on db.

