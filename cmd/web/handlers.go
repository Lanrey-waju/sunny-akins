package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/Lanrey-waju/sunny-akinns/internal/database"
)

type createContactForm struct {
	Name        string
	Email       string
	Message     string
	FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	form := createContactForm{}
	data := app.newTemplateData(r)
	data.Form = form

	app.render(w, http.StatusOK, "index.html", data)
}

func (app *application) contactMe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := createContactForm{
		Name:        r.PostForm.Get("name"),
		Email:       r.PostForm.Get("email"),
		Message:     r.PostForm.Get("content"),
		FieldErrors: map[string]string{},
	}

	// Validate form fields
	if strings.TrimSpace(form.Name) == "" {
		form.FieldErrors["name"] = "This field cannot be blank"
	} else {
		form.FieldErrors["name"] = "This field cannot be more than 50 characters long"
	}

	if strings.TrimSpace(form.Message) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	// if there are validation errors
	if len(form.FieldErrors) > 0 {
		if accepts := r.Header.Get("Accept"); strings.Contains(accepts, "application/json") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "error",
				"errors": form.FieldErrors,
				"formData": map[string]string{
					"name":    form.Name,
					"email":   form.Email,
					"content": form.Message,
				},
			})
			return
		}

		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "index.html", data)
	}

	contact, err := app.db.CreateContact(r.Context(), database.CreateContactParams{
		ID:      uuid.New(),
		Name:    form.Name,
		Email:   form.Email,
		Message: form.Message,
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
