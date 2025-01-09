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

	ts, err := template.ParseFiles("/Users/user/Documents/coding/sunny-akins/ui/html/index.html")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
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
