package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/Lanrey-waju/sunny-akinns/internal/database"
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
	if err := r.ParseForm(); err != nil {
		app.serverError(w, err)
	}
	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	content := r.PostForm.Get("content")

	contact, err := app.db.CreateContact(r.Context(), database.CreateContactParams{
		ID:      uuid.New(),
		Name:    name,
		Email:   email,
		Message: content,
	})
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("%v created and saved successfully!", contact)

	if accepts := r.Header.Get("Accept"); strings.Contains(accepts, "application/json") {
		// Return JSON response for fetch requests
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Form submitted successfully",
		})
	} else {
		// Redirect for traditional form submissions
		http.Redirect(w, r, "/", http.StatusSeeOther) // 303 See Other
	}
}
