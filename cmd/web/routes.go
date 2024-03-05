package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Post("/news/create", dynamicMiddleware.ThenFunc(app.createPost))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(loginUserHandler(app)))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))
	mux.Post("/comment/add", dynamicMiddleware.ThenFunc(app.commentNews))
	apiHandler, err := NewAPIHandler("ui/templates/api.page.tmpl")
	if err != nil {
		app.serverError(nil, err)
	}
	mux.Get("/api", dynamicMiddleware.ThenFunc(apiHandler.HandleAPI))
	mux.Get("/news/:id", dynamicMiddleware.ThenFunc(app.showNews))
	return standardMiddleware.Then(mux)
}
