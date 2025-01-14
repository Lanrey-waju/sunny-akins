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

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("POST /contact/post", app.contactMe)
	mux.HandleFunc("GET /form", app.showForm)
	mux.HandleFunc("GET /close-flash", app.closeFlashMessage)
	mux.HandleFunc("GET /about", app.aboutPage)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
}
