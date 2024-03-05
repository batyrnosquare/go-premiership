package main

import (
	"batyrnosquare/go-premiership/pkg/models"
	"batyrnosquare/go-premiership/pkg/models/mongodb"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	if err := app.errorLog.Output(2, trace); err != nil {
		app.errorLog.Println("issue with printing error logs", err)
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}

func IsAdmin() bool {
	return models.User{}.Role == string(mongodb.Admin)

}

func (app *application) adminUsers(w http.ResponseWriter, r *http.Request) {
	user := app.isAuthenticated(r)
	if user == false {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
}
