package main

import (
	"html/template"
	"log"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"/Users/user/Documents/coding/sunny-akins/ui/html/base.html",
		"/Users/user/Documents/coding/sunny-akins/ui/html/index.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	log.Print(ts.Name())
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API working good!"))
}

func (app *application) contactMe(w http.ResponseWriter, r *http.Request) {
	log.Print("contact handler invoked")
	w.Write([]byte("Send me a message!"))
}
