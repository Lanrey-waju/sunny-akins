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
	log.Print("home handler called!")

	ts, err := template.ParseFiles("./ui/html/index.html")
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

func (app *application) sayBye(w http.ResponseWriter, r *http.Request) {
	log.Print("satbye handler called!")
	w.Write([]byte("Come back again!"))
}

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Print("healthCheck handler called!")
	w.Write([]byte("API working good!"))
}

func (app *application) contactMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusBadRequest)
	}
	log.Print("contact handler invoked")
	w.Write([]byte("Send me a message!"))
}
