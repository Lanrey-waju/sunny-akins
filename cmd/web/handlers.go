package main

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/Lanrey-waju/sunny-akins/internal/database"
)

type ContactCreateForm struct {
	Name        string
	Email       string
	Message     string
	FieldErrors map[string]string
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	form := ContactCreateForm{}
	data := app.newTemplateData(r)
	data.Form = form

	if r.Header.Get("HX-Request") == "true" {
		app.renderPartial(w, "index.html", "main", data)
	} else {
		app.render(w, http.StatusOK, "index.html", data)
	}
}

func (app *application) contactMe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := ContactCreateForm{
		Name:        r.PostForm.Get("name"),
		Email:       r.PostForm.Get("email"),
		Message:     r.PostForm.Get("content"),
		FieldErrors: map[string]string{},
	}

	// Validate form fields
	if strings.TrimSpace(form.Name) == "" {
		form.FieldErrors["name"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Name) > 50 {
		form.FieldErrors["name"] = "This field cannot be more than 50 characters long"
	}

	if strings.TrimSpace(form.Message) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.renderPartial(w, "index.html", "form", data)
		return
	}

	contact, err := app.db.CreateContact(r.Context(), database.CreateContactParams{
		ID:      uuid.New(),
		Name:    form.Name,
		Email:   form.Email,
		Message: form.Message,
	})
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("%v created and saved successfully!", contact)

	app.background(func() {
		err = app.mailer.Send("hello@example.com", "request_received.tmpl", nil)
		if err != nil {
			app.errorLog.Print(err)
		}
	})

	data := app.newTemplateData(r)
	data.Form = ContactCreateForm{}
	data.Flash = "Thank You! We will get in touch."
	app.renderPartial(w, "index.html", "flash", data)
}

func (app *application) showForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = ContactCreateForm{}
	if r.Header.Get("HX-Request") == "true" {
		app.renderPartial(w, "index.html", "form", data)
		return
	}
}

func (app *application) closeFlashMessage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = ContactCreateForm{}
	w.WriteHeader(http.StatusOK)
	app.renderPartial(w, "index.html", "close-flash", data)
}

func (app *application) aboutPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = ContactCreateForm{}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Push-Url", "/about")
		w.Header().Set("HX-Title", "About Us - Sunny Akins")
		app.renderPartial(w, "about.html", "main", data)

	} else {
		app.render(w, http.StatusOK, "about.html", data)
	}
}
