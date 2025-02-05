package main

import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/Lanrey-waju/sunny-akins/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/", fileServer)

	mux.HandleFunc("GET /ping", app.ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /", dynamic.ThenFunc(app.home))
	mux.Handle("POST /contact/post", dynamic.ThenFunc(app.contactMe))
	mux.HandleFunc("GET /form", app.showForm)
	mux.HandleFunc("GET /close-flash", app.closeFlashMessage)
	mux.Handle("GET /about", dynamic.ThenFunc(app.aboutPage))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
}
